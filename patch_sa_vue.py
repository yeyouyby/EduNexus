import re

with open('edunexus/frontend/src/views/QuantumSeating.vue', 'r') as f:
    content = f.read()

# Replace ref for conflicting students
replacement1 = """const temperature = ref(100)
const iterationCount = ref(0)
const conflicts = ref(0)
const maxIter = ref(500)
const conflictingStudents = ref<number[]>([])"""

pattern1 = r"""const temperature = ref\(100\)\nconst iterationCount = ref\(0\)\nconst conflicts = ref\(0\)\nconst maxIter = ref\(500\)"""

content = re.sub(pattern1, replacement1, content)

# Update the listener
replacement2 = """      iterationCount.value = data.iteration
      temperature.value = data.temp
      conflicts.value = data.conflicts
      seats.value = data.grid
      cols.value = data.cols
      rows.value = data.rows
      if (data.conflicting_students) {
        conflictingStudents.value = data.conflicting_students
      } else {
        conflictingStudents.value = []
      }"""

pattern2 = r"""      iterationCount\.value = data\.iteration\n      temperature\.value = data\.temp\n      conflicts\.value = data\.conflicts\n      seats\.value = data\.grid\n      cols\.value = data\.cols\n      rows\.value = data\.rows"""

content = re.sub(pattern2, replacement2, content)

# Update draw logic
replacement3 = """    // Determine state purely based on current temp (jitter effect) + conflict highlight
    let inConflict = conflictingStudents.value.includes(studentId)

    if (temperature.value > 1.0) {"""

pattern3 = r"""    // Determine state purely based on current temp \(jitter effect\) \+ conflict highlight\n    let inConflict = false\n    constraints\.value\.forEach\(constRule => \{\n        if \(constRule\.student1 === studentId \|\| constRule\.student2 === studentId\) \{\n            inConflict = true\n        \}\n    \}\)\n\n    if \(temperature\.value > 1\.0\) \{"""

content = re.sub(pattern3, replacement3, content)

with open('edunexus/frontend/src/views/QuantumSeating.vue', 'w') as f:
    f.write(content)
