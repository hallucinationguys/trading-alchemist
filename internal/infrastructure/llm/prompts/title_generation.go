package prompts

// TitleGenerationData represents the data structure for title generation templates
type TitleGenerationData struct {
	UserMessage      string
	AssistantMessage string
}

// GetTitleGenerationPrompt returns the prompt configuration for conversation title generation
func GetTitleGenerationPrompt() *Prompt {
	return &Prompt{
		Name:        "title_generation",
		Description: "Generates concise, descriptive titles for conversations based on the first message exchange",
		Version:     "1.0",
		SystemPrompt: `You are a helpful assistant that generates concise, descriptive titles for conversations.
Your task is to create a title that captures the main topic or question from the conversation.

Guidelines:
- Maximum 50 characters
- Be specific and descriptive
- Use title case (capitalize major words)
- No quotes or special formatting
- Focus on the main topic or question
- Avoid generic phrases like "Help with" or "Question about"

Examples of good titles:
- "Python Data Analysis Optimization"
- "React Component Architecture"
- "SQL Query Performance Issues"
- "Machine Learning Model Selection"
- "Docker Container Configuration"
- "API Authentication Best Practices"
- "JavaScript Async/Await Patterns"
- "Database Schema Design"

Examples of bad titles:
- "Help with Python" (too vague)
- "Question about React" (too generic)
- "I need assistance with..." (too long/wordy)
- "Can you help me..." (focuses on request, not topic)`,

		UserTemplate: `Based on this conversation, generate a concise title (maximum 50 characters):

User: {{.UserMessage}}
Assistant: {{.AssistantMessage}}

Generate only the title, no additional text or formatting:`,
	}
} 