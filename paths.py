import mediapipe as mp
import json


MP_HOLISTIC = mp.solutions.holistic
FACE_MESH_EDGES = MP_HOLISTIC.FACEMESH_TESSELATION

MP_FACE_MESH = mp.solutions.face_mesh
LEFT_EYE_EDGES = MP_FACE_MESH.FACEMESH_LEFT_EYE
RIGHT_EYE_EDGES = MP_FACE_MESH.FACEMESH_RIGHT_EYE
LIPS_EDGES = MP_FACE_MESH.FACEMESH_LIPS


def findAllFaceTriangles():
    triangles = set()
    for first_edge in FACE_MESH_EDGES:
        for second_edge in FACE_MESH_EDGES:
            for iteration in range(4):
                if iteration == 0:
                    point_second_edge = 0
                    point_first_edge = 0
                elif iteration == 1:
                    point_second_edge = 1
                    point_first_edge = 1
                elif iteration == 2:
                    point_second_edge = 0
                    point_first_edge = 1
                else:
                    point_second_edge = 1
                    point_first_edge = 0

            if second_edge[point_second_edge] == first_edge[point_first_edge]:
                for third_edge in FACE_MESH_EDGES:
                    if (third_edge[0] == second_edge[(point_second_edge+1)%2] and \
                        third_edge[1] == first_edge[(point_first_edge+1)%2]) or \
                        (third_edge[1] == second_edge[(point_second_edge+1)%2] and \
                        third_edge[0] == first_edge[(point_first_edge+1)%2]):

                        triangle = frozenset((
                            first_edge[0], 
                            first_edge[1], 
                            second_edge[0], 
                            second_edge[1], 
                            third_edge[0], 
                            third_edge[1]))
                        
                        triangles.add(triangle)

    triangles = list(triangles)
    triangles = [list(element) for element in triangles]

    return triangles

def findPath(edges):
    copy_edges = edges.copy()

    for i in range(len(copy_edges)):
        copy_edges[i] = list(copy_edges[i])

    while len(copy_edges) != 0:
        path = [copy_edges[0]]

        while len(path) != len(copy_edges):
            added = False

            for i in range(len(copy_edges)):
                if copy_edges[i][0] == path[-1][1] and copy_edges[i] not in path and copy_edges[i][::-1] not in path:
                    path.append(copy_edges[i])
                    added = True
                elif copy_edges[i][1] == path[-1][1] and copy_edges[i] not in path and copy_edges[i][::-1] not in path:
                    path.append(copy_edges[i][::-1])
                    added = True

            if added == False:
                break

        allPathPoints = [point for enge in path for point in enge]

        if all(allPathPoints.count(point) == 2 for point in allPathPoints):
            return path

        copy_edges = [edge for edge in copy_edges if edge not in path]

    # FIXME - выбросить ошибку

def findPathLeftEye():
    return findPath(list(LEFT_EYE_EDGES))

def findPathRightEye():
    return findPath(list(RIGHT_EYE_EDGES))

def findPathLips():
    return findPath(list(LIPS_EDGES))

def createJSON(array_edges, file_name):
    with open(f"{file_name}.json", "w", encoding="utf-8") as file:
            json.dump(array_edges, file, indent=4)

def createJSONFaceTrianglesPath(face_triangles):
    createJSON(face_triangles, "face_triangles_path")

def createJSONLeftEyePath(path_left_eye):
    createJSON(path_left_eye, "left_eye_path")

def createJSONRigthEyePath(path_right_eye):
    createJSON(path_right_eye, "right_eye_path")

def createJSONLips(path_lips):
    createJSON(path_lips, "lips_path")

if __name__ == "__main__":
    triangles = findAllFaceTriangles()
    createJSONFaceTrianglesPath(triangles)
    print("Number of face triangles:", len(triangles))
    print("Examples of face triangles (5 numbers):")
    for triangle in triangles[:5]:
        print(triangle, "размер:", len(triangle))
    print("Max ID:", max(max(triangle) for triangle in triangles))
    print()

    path_left_eye = findPathLeftEye()
    createJSONLeftEyePath(path_left_eye)

    path_right_eye = findPathRightEye()
    createJSONRigthEyePath(path_right_eye)

    path_lips = findPathLips()
    createJSONLips(path_lips)

    print("Number of edges in path for left eye:", len(path_left_eye))
    print("Path for left eye:", path_left_eye)
    print()
    print("Number of edges in path for right eye:", len(path_right_eye))
    print("Path for right eye:",path_right_eye)
    print()
    print("Number of edges in path for lips:", len(path_lips))
    print("Path for lips:",path_lips)