#!/usr/bin/env python3
"""
Fix step numbering issues in domain guides
"""

import os
import re
import glob

def fix_step_numbering(file_path):
    """Fix step numbering in a single guide"""
    domain = os.path.basename(os.path.dirname(file_path))

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"❌ Error reading {domain}: {e}")
        return False

    original_content = content

    # Find all step headers and their positions
    step_pattern = r'## Step (\d+):'
    steps = []
    for match in re.finditer(step_pattern, content):
        step_num = int(match.group(1))
        step_text = match.group(0)
        position = match.start()
        steps.append((step_num, step_text, position))

    if not steps:
        print(f"⚠️  No steps found in {domain}")
        return False

    print(f"📋 {domain}: Found steps {[s[0] for s in steps]}")

    # Check if we have the "Using Your Own Data" section
    using_own_data_match = re.search(r'## Step \d+: Using Your Own .+ Data', content)

    if using_own_data_match:
        # Strategy: Renumber everything after "Using Your Own Data"
        using_own_data_pos = using_own_data_match.start()

        # Find steps that come after "Using Your Own Data"
        steps_after = [(s[0], s[1], s[2]) for s in steps if s[2] > using_own_data_pos]

        # Renumber steps after "Using Your Own Data"
        for old_step_num, old_step_text, pos in reversed(steps_after):  # Reverse to avoid position shifts
            if old_step_num >= 10:  # These should be incremented
                new_step_num = old_step_num + 1
                new_step_text = f"## Step {new_step_num}:"
                # Replace the specific step header
                old_pattern = re.escape(old_step_text)
                content = re.sub(old_pattern, new_step_text, content, count=1)
                print(f"  Renumbered: {old_step_text} → {new_step_text}")

        # Ensure "Using Your Own Data" is Step 9
        content = re.sub(r'## Step \d+: Using Your Own (.+) Data', r'## Step 9: Using Your Own \1 Data', content)

    else:
        print(f"⚠️  No 'Using Your Own Data' section found in {domain}")

    # Final validation: Check for duplicates
    final_steps = re.findall(r'## Step (\d+):', content)
    final_step_nums = [int(s) for s in final_steps]

    if len(final_step_nums) != len(set(final_step_nums)):
        print(f"⚠️  Still have duplicate steps in {domain}: {final_step_nums}")
        # Try a more aggressive fix
        content = fix_duplicates_aggressively(content, domain)

    # Write the fixed content
    if content != original_content:
        try:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"✅ Fixed step numbering in {domain}")
            return True
        except Exception as e:
            print(f"❌ Error writing {domain}: {e}")
            return False
    else:
        print(f"🔄 No changes needed for {domain}")
        return True

def fix_duplicates_aggressively(content, domain):
    """Aggressively fix duplicate step numbers"""
    print(f"🔧 Applying aggressive fixes to {domain}")

    # Strategy: Find all step headers and renumber them sequentially
    step_headers = []

    # Find all step headers with their full text
    for match in re.finditer(r'## Step \d+: (.+)', content):
        step_headers.append((match.group(0), match.group(1), match.start(), match.end()))

    # Sort by position
    step_headers.sort(key=lambda x: x[2])

    # Renumber sequentially
    replacements = []
    expected_step = 1

    for i, (full_header, step_title, start_pos, end_pos) in enumerate(step_headers):
        new_header = f"## Step {expected_step}: {step_title}"
        replacements.append((full_header, new_header))
        expected_step += 1

    # Apply replacements in reverse order to preserve positions
    for old_header, new_header in reversed(replacements):
        content = content.replace(old_header, new_header, 1)
        print(f"  Fixed: {old_header} → {new_header}")

    return content

def main():
    """Fix step numbering in all domain guides"""
    base_path = "/Users/scttfrdmn/src/aws-research-wizard/docs/domain-guides"
    guide_files = glob.glob(f"{base_path}/*/getting-started.md")

    print("🔧 Fixing step numbering issues...")
    print("=" * 50)

    fixed_count = 0
    for guide_file in sorted(guide_files):
        if fix_step_numbering(guide_file):
            fixed_count += 1
        print()

    print("=" * 50)
    print(f"✅ Fixed step numbering in {fixed_count}/{len(guide_files)} guides")

if __name__ == "__main__":
    main()
