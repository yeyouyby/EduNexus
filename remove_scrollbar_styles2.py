import os
import re

directory = 'edunexus/frontend/src/views'
pattern = r'<style scoped>[\s\S]*?\.custom-scrollbar::-webkit-scrollbar[\s\S]*?</style>'

for filename in os.listdir(directory):
    if filename.endswith(".vue"):
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as file:
            content = file.read()

        new_content = re.sub(pattern, '', content)

        if content != new_content:
            with open(filepath, 'w') as file:
                file.write(new_content)
            print(f"Removed custom scrollbar block from {filename}")
