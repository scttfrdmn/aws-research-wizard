# Documentation Style Guide

## 🎯 **Core Principles**

### 1. Researcher Productivity First
Every piece of documentation must help researchers get productive faster. Ask: "Does this help someone deploy a working environment in 15 minutes?"

### 2. High School Freshman Reading Level
- **Target Flesch Reading Ease**: 80+ (9th grade)
- **Sentence Length**: Maximum 20 words
- **Paragraph Length**: Maximum 4 sentences
- **Vocabulary**: Use common words, define technical terms

### 3. Learning by Doing
- Start with concrete examples, not abstract concepts
- Every guide builds something real and useful
- Include actual research scenarios with real outcomes

## ✍️ **Writing Standards**

### Language Guidelines
- **Voice**: Active voice ("Click the button" not "The button should be clicked")
- **Tense**: Present tense ("The system creates" not "The system will create")
- **Person**: Second person ("you") throughout
- **Tone**: Friendly, encouraging, non-intimidating

### Sentence Structure
```
✅ Good: "Click Deploy to start your environment."
❌ Bad: "The deployment process can be initiated by clicking the Deploy button."

✅ Good: "This costs $5 per hour."
❌ Bad: "Billing rates are approximately $5.00 per hour of usage."

✅ Good: "Type this command to connect."
❌ Bad: "Connection can be established via execution of the following command."
```

### Technical Terms
- **Define on first use**: "SSH (secure shell) lets you connect to remote computers"
- **Use analogies**: "Think of AWS regions like different data centers around the world"
- **Avoid jargon**: "Start" instead of "instantiate", "Connect" instead of "establish connection"

## 📋 **Content Structure**

### Page Template
```markdown
# [Clear, Action-Oriented Title]

> **Time**: X minutes
> **Cost**: $X-Y
> **Level**: Beginner/Intermediate/Advanced

## What You'll Build
[Concrete outcome description]

## Before You Start
[Prerequisites checklist]

## Step-by-Step Instructions
[Numbered steps with clear actions]

## What's Next
[Logical next steps]

## Troubleshooting
[Common issues and solutions]

## Feedback
[How to improve this guide]
```

### Step Format
```markdown
## Step X: [Action-Oriented Title]

[Brief explanation of what this step accomplishes]

```bash
command-to-run
```

**What this does**: [Explanation in plain English]

**Expected result**: [What the user should see]

**⚠️ If you see an error**: [Troubleshooting guidance]

**💰 Cost note**: [Billing implications if relevant]
```

## 🖼️ **Visual Standards**

### Screenshots
- **Consistency**: Same browser theme, terminal colors, window sizes
- **Annotations**: Clear arrows, numbered callouts, highlighted areas
- **Context**: Show enough of the screen for orientation
- **Alt Text**: Descriptive alt text for accessibility

### Diagrams
- **Simple**: Use basic shapes and clear labels
- **Color**: High contrast, colorblind-friendly palette
- **Text**: Large, readable fonts
- **Purpose**: Supplement text, don't replace it

### Code Blocks
- **Syntax Highlighting**: Always specify language
- **Line Length**: Wrap at 80 characters
- **Comments**: Explain non-obvious commands
- **Copy-Pasteable**: Test that commands work exactly as written

## 👥 **Personas & Scenarios**

### Researcher Personas
Each domain guide features a specific researcher:
- **Background**: Real job title, institution type, research goals
- **Current Challenge**: Specific problem they're trying to solve
- **Success Metric**: Concrete outcome they achieve
- **Quote**: Authentic-sounding testimonial

### Example Personas
**Dr. Jennifer Martinez** - Marine Biologist at UC San Diego
- Studies coral reef resilience to climate change
- Needs to process 50GB of underwater sensor data
- Currently waits 3 days for university cluster access
- "I went from 3-day waits to 30-minute analysis"

**Alex Chen** - PhD Student in Computational Chemistry
- Researching new battery materials
- Runs molecular dynamics simulations
- Laptop takes 2 weeks for calculations that cloud does in 4 hours
- "My research moved 8x faster with proper compute resources"

## 🎓 **Pedagogical Approach**

### Progressive Disclosure
1. **Start Simple**: Single domain, default settings, one instance
2. **Add Complexity**: Multiple instances, custom configurations
3. **Advanced Use**: Multi-domain, team setups, cost optimization

### Error Prevention Strategy
- **Anticipate Problems**: Based on user testing and support tickets
- **Validate Early**: Check credentials before starting deployments
- **Clear Recovery**: Specific steps to fix common errors
- **Cost Protection**: Always explain billing implications

### Learning Reinforcement
- **Repetition**: Key concepts appear in multiple contexts
- **Practice**: Each guide includes hands-on exercises
- **Validation**: Users verify each step worked correctly
- **Extension**: Clear paths to learn more

## 📊 **Quality Assurance**

### Content Review Checklist
- [ ] Reading level under 9th grade (use [Hemingway Editor](http://hemingwayapp.com/))
- [ ] All commands tested on clean environment
- [ ] Screenshots current and properly annotated
- [ ] Cost estimates accurate and up-to-date
- [ ] Time estimates realistic for beginners
- [ ] Error scenarios covered
- [ ] Mobile-friendly formatting

### User Testing Protocol
1. **Recruit**: Find researchers unfamiliar with AWS
2. **Observe**: Watch them follow guide without help
3. **Note**: Where they get confused or stuck
4. **Time**: Measure actual completion times
5. **Survey**: Ask about clarity and confidence
6. **Iterate**: Update guide based on feedback

### Accessibility Standards
- **Screen Readers**: Proper heading hierarchy, alt text
- **Color**: Don't rely solely on color to convey information
- **Contrast**: WCAG 2.1 AA compliant contrast ratios
- **Mobile**: Readable and usable on phones/tablets
- **Cognitive**: Clear structure, consistent navigation

## 🔄 **Maintenance & Updates**

### Regular Review Schedule
- **Monthly**: Check all external links and screenshots
- **Quarterly**: Verify cost estimates and time requirements
- **Semi-Annually**: Full user testing and content refresh
- **Annually**: Complete style guide review

### Version Control
- Track major changes to content
- Note when screenshots were last updated
- Document cost estimate revision dates
- Maintain changelog for significant updates

### Community Feedback Integration
- **GitHub Issues**: Template for documentation problems
- **User Surveys**: Embedded feedback forms
- **Analytics**: Track where users drop off
- **Support Tickets**: Identify common confusion points

This style guide ensures consistent, accessible, and effective documentation that helps researchers succeed quickly while identifying areas for improvement.
