import cv2 as cv
import mediapipe as mp
import numpy as np
import json
import subprocess
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os


BASE_DIR = os.path.dirname(os.path.abspath(__file__))


def get_path(file_name):
    with open(file_name, "r", encoding="utf-8") as file:
        return json.load(file)


faceTrianglesPath = get_path(os.path.join(BASE_DIR, "json", "face_triangles_path.json"))
leftEyePath = get_path(os.path.join(BASE_DIR, "json", "left_eye_path.json"))
rigthEyePath = get_path(os.path.join(BASE_DIR, "json", "right_eye_path.json"))
lipsPath = get_path(os.path.join(BASE_DIR, "json", "lips_path.json"))


app = FastAPI()


class VideoDTO(BaseModel):
    file_name: str
    file_path: str


class VideoDomain(BaseModel):
    file_name: str
    file_path: str
    file_size: int


@app.patch("/videos")
def videoBlur(videoDTO: VideoDTO) -> VideoDomain:
    cap = cv.VideoCapture(videoDTO.file_path)
    if not cap.isOpened():
        raise HTTPException(400, "Cannot open file")

    fps = cap.get(cv.CAP_PROP_FPS)
    width = int(cap.get(cv.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv.CAP_PROP_FRAME_HEIGHT))

    folder = os.getenv("BLURRED_VIDEOS_FOLDER", "../blurred_videos_folder")
    os.makedirs(folder, exist_ok=True)

    name = videoDTO.file_name
    if not name.endswith(".mp4"):
        name += ".mp4"
    newFileName = f"blurred_{name}"
    newFilePath = os.path.join(folder, newFileName)

    # 1. Пишем во временный файл (mp4v)
    tmpPath = newFilePath + ".tmp.mp4"
    fourcc = cv.VideoWriter_fourcc(*"mp4v")
    out = cv.VideoWriter(tmpPath, fourcc, fps, (width, height))

    if not out.isOpened():
        raise HTTPException(500, "Cannot open VideoWriter")

    MP_FACE_MESH = mp.solutions.face_mesh

    with MP_FACE_MESH.FaceMesh(
        max_num_faces=15,
        min_detection_confidence=0.3,
        min_tracking_confidence=0.3,
    ) as face_mesh:
        while cap.isOpened():
            ret, frame = cap.read()
            if not ret:
                break

            imageRGB = cv.cvtColor(frame, cv.COLOR_BGR2RGB)
            results = face_mesh.process(imageRGB)
            imageBGR = cv.cvtColor(imageRGB, cv.COLOR_RGB2BGR)

            if results.multi_face_landmarks:
                h, w = imageBGR.shape[:2]

                mask = np.zeros(imageBGR.shape[:2], dtype=np.uint8)
                blurred = cv.GaussianBlur(imageBGR, (51, 51), 0)

                for face_landmarks in results.multi_face_landmarks:
                    points = []

                    for landmark in face_landmarks.landmark:
                        current_x = int(landmark.x * w)
                        current_y = int(landmark.y * h)
                        points.append([current_x, current_y])

                    for i in range(len(faceTrianglesPath)):
                        np_polygon = np.array([
                            points[faceTrianglesPath[i][0]],
                            points[faceTrianglesPath[i][1]],
                            points[faceTrianglesPath[i][2]],
                        ], dtype=np.int32)
                        cv.fillPoly(mask, [np_polygon], 255)

                    np_polygon = np.array(
                        [points[element] for pair in leftEyePath for element in pair],
                        dtype=np.int32,
                    )
                    cv.fillPoly(mask, [np_polygon], 255)

                    np_polygon = np.array(
                        [points[element] for pair in rigthEyePath for element in pair],
                        dtype=np.int32,
                    )
                    cv.fillPoly(mask, [np_polygon], 255)

                    np_polygon = np.array(
                        [points[element] for pair in lipsPath for element in pair],
                        dtype=np.int32,
                    )
                    cv.fillPoly(mask, [np_polygon], 255)

                imageBGR[mask == 255] = blurred[mask == 255]

            out.write(imageBGR)

    out.release()
    cap.release()

    # 2. Конвертируем mp4v
    try:
        subprocess.run([
            "ffmpeg", "-y",
            "-i", tmpPath,
            "-c:v", "libx264",
            "-preset", "fast",
            "-crf", "23",
            "-pix_fmt", "yuv420p",
            "-movflags", "+faststart",
            newFilePath,
        ], check=True, capture_output=True)
    except subprocess.CalledProcessError as e:
        raise HTTPException(500, f"ffmpeg error: {e.stderr.decode()}")
    finally:
        if os.path.exists(tmpPath):
            os.remove(tmpPath)

    newFileSize = os.path.getsize(newFilePath)

    return VideoDomain(
        file_name=newFileName,
        file_path=newFilePath,
        file_size=newFileSize,
    )