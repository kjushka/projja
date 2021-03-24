import math


def dotry(i: int) -> bool:
    if (vx[i]):
        return False
    vx[i] = True
    for j in range(n):
        if (a[i][j] - maxrow[i] - mincol[j] == 0):
            vy[j] = True
    for j in range(n):
        if (a[i][j] - maxrow[i] - mincol[j] == 0 & yx[j] == -1):
            xy[i] = j
            yx[j] = i
            return True
    for j in range(n):
        if (a[i][j] - maxrow[i] - mincol[j] == 0 & dotry(yx[j])):
            xy[i] = j
            yx[j] = i
            return True
    return False


def main():

    for i in range(n):
        maxrow.append(0)
        mincol.append(0)
        xy.append(-1)
        yx.append(-1)
    for i in range(n):
        for j in range(n):
            maxrow[i] = max(maxrow[i], a[i][j])
    for c in range(0, n):
        for j in range(n):
            vx.append(0)
            vy.append(0)
        k = 0
        for i in range(n):
            if (xy[i] == -1 & dotry(i)):
                k += 1
        c += k
        if (k == 0):
            z = math.inf
            for i in range(n):
                if (vx[i]):
                    for j in range(n):
                        if (not vy[j]):
                            z = min(z, maxrow[i] + mincol[j] - a[i][j])
            for i in range(n):
                if (vx[i]):
                    maxrow[i] -= z
                if (vy[i]):
                    mincol[i] += z

    ans = 0
    for i in range(n):
        ans += a[i][xy[i]]
    print('{}\n'.format(ans))
    for i in range(n):
        print("{} ".format(xy[i] + 1))

n = 5
a = list()
a.append([1, 2, 3, 4, 0])
a.append([4, 0, 1, 2, 3])
a.append([2, 4, 3, 0, 2])
a.append([2, 1, 1, 3, 4])
a.append([0, 3, 4, 0, 0])
xy = list()
yx = list()
vx = list()
vy = list()
maxrow = list()
mincol = list()

if (__name__ == "__main__"):
    main()