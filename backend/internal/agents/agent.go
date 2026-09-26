package agents

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AgentType represents the type of agent
type AgentType string

const (
	AgentIntelCollector AgentType = "intel_collector"
	AgentScriptWriter   AgentType = "script_writer"
	AgentXPostGenerator AgentType = "x_post_generator"
	AgentThumbnailForge AgentType = "thumbnail_forge"
	AgentBlogWriter     AgentType = "blog_writer"
	AgentOutreachEngine AgentType = "outreach_engine"
	AgentNewsletter     AgentType = "newsletter"
	AgentClipFactory    AgentType = "clip_factory"
	AgentPerformanceEval AgentType = "performance_eval"
	AgentWeeklyDigest   AgentType = "weekly_digest"
)

// Agent represents a content generation agent
type Agent struct {
	ID          string                 `json:"id"`
	Type        AgentType              `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      AgentStatus            `json:"status"`
	Prompt      string                 `json:"prompt"`
	Model       string                 `json:"model"`
	Provider    string                 `json:"provider"`
	Config      map[string]interface{} `json:"config"`
	LastRun     *time.Time             `json:"lastRun,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}

// AgentStatus represents the agent's current state
type AgentStatus string

const (
	AgentStatusIdle      AgentStatus = "idle"
	AgentStatusRunning   AgentStatus = "running"
	AgentStatusDisabled  AgentStatus = "disabled"
	AgentStatusError     AgentStatus = "error"
)

// Result represents the output of an agent's work
type Result struct {
	AgentID     string                 `json:"agentId"`
	AgentType   AgentType              `json:"agentType"`
	ContentID   string                 `json:"contentId"`
	ContentType string                 `json:"contentType"`
	Title       string                 `json:"title"`
	Body        string                 `json:"body"`
	Metadata    map[string]interface{} `json:"metadata"`
	Error       string                 `json:"error,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Timestamp   time.Time              `json:"timestamp"`
}

// LLMClient interface for agents to use
type LLMClient interface {
	Complete(ctx context.Context, prompt string, opts RequestOptions) (string, error)
}

// RequestOptions for LLM calls
type RequestOptions struct {
	Model       string
	MaxTokens   int
	Temperature float64
	System      string
	Messages    []Message
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// New creates a new agent
func New(agentType AgentType, config map[string]interface{}) *Agent {
	return &Agent{
		ID:          generateID(string(agentType)),
		Type:        agentType,
		Name:        getAgentName(agentType),
		Description: getAgentDescription(agentType),
		Status:      AgentStatusIdle,
		Prompt:      getAgentPrompt(agentType),
		Config:      config,
		CreatedAt:   time.Now(),
	}
}

// Execute runs the agent
func (a *Agent) Execute(ctx context.Context, client LLMClient, input map[string]interface{}) (*Result, error) {
	startTime := time.Now()
	
	a.Status = AgentStatusRunning
	now := time.Now()
	a.LastRun = &now
	
	// Build the prompt with context
	prompt := a.buildPrompt(input)
	
	// Call LLM
	response, err := client.Complete(ctx, prompt, RequestOptions{
		Model:       a.Model,
		Temperature: 0.7,
	})
	
	if err != nil {
		a.Status = AgentStatusError
		return &Result{
			AgentID:   a.ID,
			AgentType: a.Type,
			Error:     err.Error(),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, err
	}
	
	a.Status = AgentStatusIdle
	
	// Parse and structure the response
	result := &Result{
		AgentID:   a.ID,
		AgentType: a.Type,
		Body:      response,
		Duration:  time.Since(startTime),
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	
	// Extract title from response if applicable
	if title := extractTitle(response); title != "" {
		result.Title = title
	}
	
	return result, nil
}

func (a *Agent) buildPrompt(input map[string]interface{}) string {
	prompt := a.Prompt
	
	// Add context from input
	for key, value := range input {
		placeholder := fmt.Sprintf("{{%s}}", key)
		prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("%v", value))
	}
	
	return prompt
}

func generateID(agentType string) string {
	return fmt.Sprintf("%s_%d", agentType, time.Now().Unix())
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			return line
		}
	}
	return ""
}

func getAgentName(agentType AgentType) string {
	names := map[AgentType]string{
		AgentIntelCollector:  "Intel Collector",
		AgentScriptWriter:    "Script Writer",
		AgentXPostGenerator:  "X Post Generator",
		AgentThumbnailForge:  "Thumbnail Forge",
		AgentBlogWriter:      "Blog Writer",
		AgentOutreachEngine:  "Outreach Engine",
		AgentNewsletter:      "Newsletter Composer",
		AgentClipFactory:     "Clip Factory",
		AgentPerformanceEval: "Performance Evaluator",
		AgentWeeklyDigest:    "Weekly Digest",
	}
	
	if name, ok := names[agentType]; ok {
		return name
	}
	return string(agentType)
}

func getAgentDescription(agentType AgentType) string {
	descriptions := map[AgentType]string{
		AgentIntelCollector:  "Collects and analyzes trending topics and competitor content",
		AgentScriptWriter:    "Writes engaging video scripts optimized for retention",
		AgentXPostGenerator:  "Generates viral X posts and threads",
		AgentThumbnailForge:  "Creates thumbnail concepts and descriptions",
		AgentBlogWriter:      "Writes SEO-optimized blog posts",
		AgentOutreachEngine:  "Composes outreach messages for collaborations",
		AgentNewsletter:      "Composes engaging newsletters",
		AgentClipFactory:     "Creates clip scripts for short-form content",
		AgentPerformanceEval: "Analyzes content performance and suggests improvements",
		AgentWeeklyDigest:    "Compiles weekly content summaries",
	}
	
	if desc, ok := descriptions[agentType]; ok {
		return desc
	}
	return ""
}

func getAgentPrompt(agentType AgentType) string {
	prompts := map[AgentType]string{
		AgentIntelCollector: `You are an expert content intelligence researcher. Your task is to:
1. Identify trending topics in {{niche}}
2. Analyze competitor content from {{competitors}}
3. Find content gaps and opportunities
4. Extract viral patterns and hooks

Focus on data from YouTube, X, and relevant industry sources.

Output format:
## Trending Topics
- [Topic]: [Why it's trending] - [Viral potential: X/10]

## Competitor Analysis
- [Channel]: [Latest video] - [Views: X] - [Pattern: Y]

## Content Gaps
- [Opportunity]: [Why it's missing] - [Your angle]

## Viral Patterns
- [Pattern]: [Example] - [How to use it]`,

		AgentScriptWriter: `You are an expert YouTube scriptwriter with 10+ years experience. Your task is to write a compelling video script.

Niche: {{niche}}
Topic: {{topic}}
Target Length: {{length}} seconds
Target Audience: {{audience}}
Tone: {{tone}}
Style: {{style}}

Structure:
1. Hook (0-30s) - Pattern interrupt, curiosity gap
2. Promise (30-45s) - What they'll learn
3. Credibility (45-60s) - Why trust you
4. Content Body (60-{{length}}s) - Value delivery with pattern interrupts every 30-45s
5. CTA (final 30s) - Subscribe, comment, watch next

Include:
- [PAUSE] markers for pacing
- [B-ROLL: description] for visual cues
- [ZOOM] for emphasis moments
- Mid-roll hook at 30% and 60% mark

Output the complete script with timing markers.`,

		AgentXPostGenerator: `You are a viral X (Twitter) content expert. Generate high-engagement posts.

Niche: {{niche}}
Topic: {{topic}}
Style: {{style}} (hook-first / educational / contrarian / story)
Posts needed: {{count}}

For each post:
1. Hook line (must stop the scroll)
2. Body (3-5 punchy lines or bullet points)
3. Engagement trigger (question or controversial take)
4. Thread hook (if applicable)

Examples of viral structures:
- "Most people think X. Here's the truth:"
- "I spent $X to learn Y. Here's what nobody told me:"
- "Stop doing X. Do this instead:"
- "In 2024, [niche] changed forever."

Generate {{count}} distinct posts with different angles on this topic.`,

		AgentThumbnailForge: `You are a YouTube thumbnail expert. Create high-CTR thumbnail concepts.

Video Title: {{title}}
Video Topic: {{topic}}
Target Audience: {{audience}}
Thumbnail Style: {{style}} (high-contrast / reaction / text-heavy / minimal)
Variants needed: {{count}}

For each thumbnail concept:
1. Layout description (subject positioning, background)
2. Text overlay (3 words max, font style, position)
3. Color scheme (high contrast colors)
4. Expression/Action (for face thumbnails)
5. CTR prediction (based on current trends)

Create {{count}} distinct concepts using different psychological triggers:
- Curiosity gap
- Shock/surprise
- FOMO
- Transformation
- Exclusivity

Include technical specs: 1280x720, faces should be 30-40% of frame.`,

		AgentBlogWriter: `You are an expert SEO blog writer. Create comprehensive, rankable articles.

Niche: {{niche}}
Topic: {{topic}}
Target Keyword: {{keyword}}
Word Count: {{count}}
Tone: {{tone}}

Structure:
1. Title (60 chars, keyword at start)
2. Meta description (155 chars, keyword + CTA)
3. Introduction (hook + problem + promise)
4. H2 sections (6-8 sections, keyword variations)
5. H3 subsections for depth
6. Conclusion with CTA
7. FAQ section (People Also Ask optimization)

SEO requirements:
- Keyword density: 1-2%
- Internal linking opportunities
- External authority links (2-3)
- Image alt text suggestions
- Schema markup recommendations

Write the complete article following EEAT principles.`,

		AgentOutreachEngine: `You are an expert outreach copywriter. Write personalized collaboration messages.

Your Channel: {{channel}}
Their Channel: {{target_channel}}
Their Content: {{target_content}}
Outreach Type: {{type}} (collab / guest / shoutout / sponsorship)

Write a message that:
1. Opens with specific genuine compliment (mention exact video/timestamp)
2. Establishes credibility (your relevant achievement)
3. Proposes win-win collaboration
4. Makes it easy to say yes (low friction ask)
5. Includes social proof

Keep it under 150 words. Avoid spam patterns. Personalize deeply.`,

		AgentNewsletter: `You are an expert newsletter writer. Create engaging email content.

Niche: {{niche}}
Topic: {{topic}}
Subscriber Segment: {{segment}}
Tone: {{tone}}

Structure:
1. Subject line (40-50 chars, emoji optional, curiosity or benefit)
2. Preview text (90-100 chars)
3. Opening hook (personal story or bold statement)
4. Main content (scannable with bullet points)
5. Call to action (single focus)
6. P.S. line (bonus value or engagement hook)

Include:
- Personal touch
- Exclusive insight
- Clear value proposition
- Single primary CTA

Write the complete newsletter ready to send.`,

		AgentClipFactory: `You are a short-form content expert. Create viral clip scripts.

Source Video: {{source_video}}
Video Topic: {{topic}}
Platform: {{platform}} (tiktok / reels / shorts)
Target Duration: {{duration}} seconds

Identify 3-5 clip-worthy moments:
1. Timestamp (start-end)
2. Hook rewrite (first 3 seconds must be explosive)
3. Caption text (max 3 lines, engaging font)
4. Hashtag strategy (3-5 tags)
5. Background music vibe
6. Expected performance (views prediction)

Focus on:
- Pattern interrupts
- Controversial takes
- Transformation moments
- Educational nuggets
- Emotional peaks

Provide timestamped cut points and rewrite hooks for maximum retention.`,

		AgentPerformanceEval: `You are a content performance analyst. Review and suggest improvements.

Channel: {{channel}}
Period: {{period}}
Top Content: {{top_content}}
Underperformers: {{underperformers}}

Analysis:
1. Metrics Overview (views, CTR, retention, engagement rates)
2. What's Working (patterns in top performers)
3. What's Not Working (common issues in underperformers)
4. Content Gaps (what you're not covering)
5. Audience Insights (demographics, behavior)
6. Competitor Moves (what's working for them)

Recommendations:
1. Immediate Actions (this week)
2. Content Strategy Shifts (next month)
3. Format Experiments to Try
4. Publishing Schedule Optimization
5. Title/Thumbnail Optimization

Be data-driven and specific.`,

		AgentWeeklyDigest: `You are a content strategist. Create a weekly performance summary.

Period: {{period}}
Content Published: {{published_count}}
Total Views: {{total_views}}
New Subscribers: {{new_subscribers}}
Top Performing: {{top_performer}}

Structure:
1. Week in Numbers (key metrics vs last week)
2. Content Highlights (each piece performance)
3. Audience Growth (subscriber trends)
4. Engagement Winners (comments, shares analysis)
5. Lessons Learned
6. Next Week Focus (3 priorities)
7. Content Calendar Preview

Format as a clean, executive summary style report suitable for team review.`,
	}
	
	if prompt, ok := prompts[agentType]; ok {
		return prompt
	}
	return "Generate content for: " + string(agentType)
}
