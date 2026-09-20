import mediapipe as mp
import json


MP_HOLISTIC = mp.solutions.holistic
FACEMESH = MP_HOLISTIC.FACEMESH_TESSELATION


def findAllTriangles():
    triangles = set()
    for first_edge in FACEMESH:
        for second_edge in FACEMESH:
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
                for third_edge in FACEMESH:
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


def createFileTriangles(triangles):
    with open("triangles.json", "w", encoding="utf-8") as file:
        json.dump(triangles, file, indent=4)


if __name__ == "__main__":
    triangles = findAllTriangles()

    createFileTriangles(triangles)

    print("Количество треугольников:", len(triangles), "\n")

    for triangle in triangles[:20]:
        print(triangle, "размер:", len(triangle))

    print()

    print("Максимальный ID:", max(max(triangle) for triangle in triangles))