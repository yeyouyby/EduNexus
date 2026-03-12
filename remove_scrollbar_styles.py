import os
import re

directory = 'edunexus/frontend/src/views'
pattern = r'<style scoped>[\s\S]*?\.custom-scrollbar::-webkit-scrollbar \{[\s\S]*?\}[\s\S]*?</style>'

for filename in os.listdir(directory):
    if filename.endswith(".vue"):
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as file:
            content = file.read()

        # Remove the <style scoped> block entirely if it only contains the scrollbar styles
        # or remove just the scrollbar styles if there are others.
        # Let's check what's actually in these files first.

        # A simple replacement that works for the exact block provided in QuantumSeating.vue
        exact_block_pattern = r'<style scoped>\n\.custom-scrollbar::-webkit-scrollbar \{\n  width: 6px;\n\}\n\.custom-scrollbar::-webkit-scrollbar-track \{\n  background: rgba\(0, 0, 0, 0\.3\);\n\}\n\.custom-scrollbar::-webkit-scrollbar-thumb \{\n  background: rgba\(0, 255, 204, 0\.2\);\n  border-radius: 3px;\n\}\n\.custom-scrollbar::-webkit-scrollbar-thumb:hover \{\n  background: rgba\(0, 255, 204, 0\.5\);\n\}\n</style>\n?'

        new_content = re.sub(exact_block_pattern, '', content)

        if content != new_content:
            with open(filepath, 'w') as file:
                file.write(new_content)
            print(f"Removed scrollbar styles from {filename}")
        else:
            # Let's try a more relaxed pattern
            relaxed_pattern = r'<style scoped>\s*\.custom-scrollbar::-webkit-scrollbar \{\s*width: 6px;\s*\}\s*\.custom-scrollbar::-webkit-scrollbar-track \{\s*background: rgba\(0, 0, 0, 0\.3\);\s*\}\s*\.custom-scrollbar::-webkit-scrollbar-thumb \{\s*background: rgba\(0, 255, 204, 0\.2\);\s*border-radius: 3px;\s*\}\s*\.custom-scrollbar::-webkit-scrollbar-thumb:hover \{\s*background: rgba\(0, 255, 204, 0\.5\);\s*\}\s*</style>\n?'
            new_content = re.sub(relaxed_pattern, '', content)
            if content != new_content:
                with open(filepath, 'w') as file:
                    file.write(new_content)
                print(f"Removed scrollbar styles (relaxed match) from {filename}")
            else:
                 print(f"No match found in {filename}")
