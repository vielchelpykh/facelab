import cv2 as cv
import mediapipe as mp
import numpy as np
import json
from fastapi import FastAPI, status, HTTPException
from pydantic import BaseModel
import os


def getPath(file_name):
            with open(file_name, "r", encoding="utf-8") as file:
                edges = json.load(file)
                return edges

def getFaceTrianglesPath():
    return getPath("./json/face_triangles_path.json")

def getLeftEyePath():
    return getPath("./json/left_eye_path.json")

def getRigthEyePath():
    return getPath("./json/right_eye_path.json")

def getLipsPath():
    return getPath("./json/lips_path.json")


faceTrianglesPath = getFaceTrianglesPath()
leftEyePath = getLeftEyePath()
rigthEyePath = getRigthEyePath()
lipsPath = getLipsPath()


app = FastAPI()


class VideoDTO(BaseModel):
    fileName: str
    filePath: str

class VideoDomain(BaseModel):
    fileName: str
    filePath: str
    fileSize: int


@app.patch("/videos")
def videoBlur(videoDTO: VideoDTO) -> VideoDomain:
    cap = cv.VideoCapture(videoDTO.filePath)
    if not cap.isOpened():
         raise HTTPException(400, "Cannot open file")

    fps = cap.get(cv.CAP_PROP_FPS)
    width = int(cap.get(cv.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv.CAP_PROP_FRAME_HEIGHT))


    # folder = os.getenv("BLURRED_VIDEOS_FOLDER")
    folder = "../blurred_videos_folder"
    os.makedirs(folder, exist_ok=True)

    newFileName = f"blurred_{videoDTO.fileName}.mp4"
    newFilePath = os.path.join(folder, newFileName)

    fourcc = cv.VideoWriter_fourcc(*"mp4v")
    out = cv.VideoWriter(newFilePath, fourcc, fps, (width, height))
    
    MP_HOLISTIC = mp.solutions.holistic

    with MP_HOLISTIC.Holistic(min_detection_confidence=0.5, min_tracking_confidence=0.5) as holistic:
        while cap.isOpened():
            ret, frame = cap.read()
            if not ret:
                break

            imageRGB = cv.cvtColor(frame, cv.COLOR_BGR2RGB)

            results = holistic.process(imageRGB)

            imageBGR = cv.cvtColor(imageRGB, cv.COLOR_RGB2BGR)

            
            if results.face_landmarks:
                points = []
                for id, landmark in enumerate(results.face_landmarks.landmark):
                    height, width = imageBGR.shape[:2]
                    current_x = int(landmark.x * width)
                    current_y = int(landmark.y * height)
                    points.append([current_x, current_y])

            blurred = cv.GaussianBlur(imageBGR, (51, 51), 0)
            mask = np.zeros(imageBGR.shape[:2], dtype=np.uint8)

            for i in range(len(faceTrianglesPath)):
                np_polygon = np.array([
                    points[faceTrianglesPath[i][0]], 
                    points[faceTrianglesPath[i][1]],
                    points[faceTrianglesPath[i][2]]],
                    dtype=np.int32,
                    )       
                cv.fillPoly(mask, [np_polygon], 255)

            np_polygon = np.array([points[element] for pair in leftEyePath for element in pair], dtype=np.int32)   
            cv.fillPoly(mask, [np_polygon], 255)

            np_polygon = np.array([points[element] for pair in rigthEyePath for element in pair], dtype=np.int32)   
            cv.fillPoly(mask, [np_polygon], 255)

            np_polygon = np.array([points[element] for pair in lipsPath for element in pair], dtype=np.int32)   
            cv.fillPoly(mask, [np_polygon], 255)
                
            imageBGR[mask == 255] = blurred[mask == 255]

            out.write(imageBGR)
            
    out.release()
    cap.release()

    newFileSize = os.path.getsize(newFilePath)

    return VideoDomain (
        fileName=newFileName,
        filePath=newFilePath,
        fileSize=newFileSize,
    )
