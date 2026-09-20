import cv2 as cv
import mediapipe as mp
import numpy as np
import json


VIDEO = cv.VideoCapture(0)
MP_HOLISTIC = mp.solutions.holistic

def getTriangles():
    with open("triangles.json", "r", encoding="utf-8") as file:
        triangles = json.load(file)
        return triangles


with MP_HOLISTIC.Holistic(min_detection_confidence=0.5, min_tracking_confidence=0.5) as holistic:
    triangles = getTriangles()

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

        overlay = imageBGR.copy()

        for i in range (len(triangles)):
            np_polygon = np.array([
                points[triangles[i][0]], 
                points[triangles[i][1]],
                points[triangles[i][2]]],
                dtype=np.int32,
                )       
            cv.fillPoly(overlay, [np_polygon], (0, 0, 255))
            

        cv.addWeighted(overlay, 0.4, imageBGR, 0.6, 0, imageBGR)

        cv.imshow("Camera", imageBGR)

        if (cv.waitKey(10) & 0xFF) == ord('q'):
            break

VIDEO.release()
cv.destroyAllWindows()