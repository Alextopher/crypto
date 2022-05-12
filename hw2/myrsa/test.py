# Returns d, x, y such that d = x * a + y * b
def xgcd(a, b):
    saved_a = a
    x1, y1, x2, y2 = 1, 0, 0, 1

    while b != 0:
        q = int(a / b)
        r = int(a % b)
        print(a, b, q, r, x1, y1, x2, y2)
        a, b = b, r
        x1, x2 = x2, x1 - q * x2
        y1, y2 = y2, y1 - q * y2

    if x1 < 0:
        x1 += saved_a
    
    if y1 < 0:
        y1 += saved_a

    return a, x1, y1

print(xgcd(1054609920, 65537))