import cv2 as cv
import mediapipe as mp
import numpy as np
import json


VIDEO = cv.VideoCapture(0)
MP_HOLISTIC = mp.solutions.holistic

def getPath(file_name):
    with open(file_name, "r", encoding="utf-8") as file:
        edges = json.load(file)
        return edges

def getFaceTrianglesPath():
    return getPath("face_triangles_path.json")

def getLeftEyePath():
    return getPath("left_eye_path.json")

def getRigthEyePath():
    return getPath("right_eye_path.json")

def getLipsPath():
    return getPath("lips_path.json")


with MP_HOLISTIC.Holistic(min_detection_confidence=0.5, min_tracking_confidence=0.5) as holistic:
    faceTrianglesPath = getFaceTrianglesPath()
    leftEyePath = getLeftEyePath()
    rigthEyePath = getRigthEyePath()
    lipsPath = getLipsPath()

    while VIDEO.isOpened():
        ret, frame = VIDEO.read()

        imageRGB = cv.cvtColor(frame, cv.COLOR_BGR2RGB)

        results = holistic.process(imageRGB)

        imageBGR = cv.cvtColor(imageRGB, cv.COLOR_RGB2BGR)

        points = []
        if results.face_landmarks:
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

        cv.imshow("Camera", imageBGR)

        if (cv.waitKey(10) & 0xFF) == ord('q'):
            break

VIDEO.release()
cv.destroyAllWindows()