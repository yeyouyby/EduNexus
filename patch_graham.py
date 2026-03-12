import re

with open('edunexus/graham_scan.go', 'r') as f:
    content = f.read()

replacement = """		// Only keep the furthest point for collinear points (points that share the same angle from p0)
		m := 0
		for i := 0; i < len(pts); {
			j := i
			// Advance j to the furthest point among those collinear with pts[i] from p0
			for j+1 < len(pts) && crossProduct(p0, pts[j], pts[j+1]) == 0 {
				j++
			}
			// Keep only the furthest point in this collinear run
			pts[m] = pts[j]
			m++
			i = j + 1
		}
		pts = pts[:m]"""

pattern = r"""\t\t// Only keep the furthest point for collinear points\n\t\tm := 1\n\t\tfor i := 1; i < len\(pts\); i\+\+ \{\n\t\t\tfor i < len\(pts\)-1 && crossProduct\(p0, pts\[i\], pts\[i\+1\]\) == 0 \{\n\t\t\t\ti\+\+\n\t\t\t\}\n\t\t\tpts\[m-1\] = pts\[i\]\n\t\t\tm\+\+\n\t\t\}\n\t\tpts = pts\[:m-1\]"""

new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

with open('edunexus/graham_scan.go', 'w') as f:
    f.write(new_content)
