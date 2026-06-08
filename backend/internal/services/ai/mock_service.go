package ai

import (
	"fmt"
	"strings"
	"time"

	"github.com/research-paper-analyzer/backend/internal/models"
)

// MockAIService provides realistic AI responses without requiring AWS credentials.
// Use this for local development and testing.
type MockAIService struct{}

// NewMockAIService creates a new mock AI service instance.
func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

// GenerateSummary returns a realistic mock analysis of the provided text.
// It adapts the response slightly based on the input text length and content.
func (m *MockAIService) GenerateSummary(text string) (*AnalysisResult, error) {
	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	// Extract some words from the text for a slightly personalized response
	words := strings.Fields(text)
	wordCount := len(words)

	// Pick topic hints from the text for a more contextual response
	topicHint := "the research topic"
	if wordCount > 10 {
		// Use first few meaningful words as a topic hint
		topicWords := words[:min(15, len(words))]
		topicHint = strings.Join(topicWords, " ")
	}

	return &AnalysisResult{
		Summary: fmt.Sprintf(
			"This research paper presents a comprehensive examination of %s. "+
				"The study spans approximately %d words of detailed analysis, methodology, and findings. "+
				"The authors provide a thorough literature review situating their work within the broader "+
				"academic context, followed by a rigorous methodological framework. The paper presents "+
				"compelling evidence through both quantitative and qualitative analyses, demonstrating "+
				"significant contributions to the field. The conclusions drawn are well-supported by the "+
				"data presented, and the authors identify important implications for both theory and practice. "+
				"Overall, this work represents a valuable addition to the existing body of knowledge and "+
				"opens up several promising avenues for future investigation.",
			topicHint, wordCount,
		),
		KeyFindings: `• The study identifies significant patterns and trends that have not been previously documented in the existing literature.
• A novel analytical framework is proposed that demonstrates improved performance compared to baseline approaches.
• Statistical analysis reveals strong correlations between the primary variables under investigation (p < 0.05).
• The results suggest practical implications for industry practitioners and policymakers in the field.
• Cross-validation experiments confirm the robustness and generalizability of the proposed approach across multiple datasets.`,
		Methodology: "The research employs a mixed-methods approach combining quantitative analysis with " +
			"qualitative insights. Data was collected through a systematic process involving multiple sources " +
			"to ensure triangulation and validity. The quantitative component utilizes advanced statistical " +
			"techniques including regression analysis, hypothesis testing, and cross-validation. The qualitative " +
			"component involves thematic analysis of expert interviews and case studies. The study follows " +
			"established research protocols and ethical guidelines, with appropriate measures taken to ensure " +
			"data quality and reliability throughout the research process.",
		Limitations: `• The sample size, while adequate for the chosen methodology, may limit the generalizability of findings to broader populations.
• The study is cross-sectional in nature, preventing causal inferences about the relationships between variables.
• Self-reported data in the qualitative component may be subject to recall bias and social desirability effects.
• The study was conducted within a specific geographic and temporal context, which may affect external validity.`,
		FutureScope: `• Longitudinal studies should be conducted to examine how the identified patterns evolve over time and establish causal relationships.
• Replication of this study across different geographic regions and cultural contexts would strengthen the external validity of findings.
• Integration of emerging technologies such as machine learning and big data analytics could enhance the analytical framework.
• Collaborative interdisciplinary research combining insights from related fields would provide a more holistic understanding of the phenomena under study.`,
		Keywords: "research analysis, methodology, data science, systematic review, statistical analysis, machine learning, evidence-based, peer-reviewed",
	}, nil
}

// GenerateQuiz returns realistic mock quiz questions based on the paper content.
func (m *MockAIService) GenerateQuiz(text string, numQuestions int) ([]models.QuizQuestionAI, error) {
	// Simulate processing time
	time.Sleep(300 * time.Millisecond)

	// Ensure we don't exceed 10 questions
	if numQuestions <= 0 {
		numQuestions = 5
	}
	if numQuestions > 10 {
		numQuestions = 10
	}

	// Pool of realistic quiz questions
	questionPool := []models.QuizQuestionAI{
		{
			Question:      "What is the primary research methodology used in this study?",
			OptionA:       "Purely qualitative case study analysis",
			OptionB:       "Mixed-methods approach combining quantitative and qualitative analysis",
			OptionC:       "Randomized controlled experiment",
			OptionD:       "Historical archival research",
			CorrectAnswer: "B",
			Explanation:   "The paper employs a mixed-methods approach, combining quantitative statistical analysis with qualitative insights from expert interviews and case studies to ensure comprehensive coverage of the research questions.",
		},
		{
			Question:      "According to the paper, which factor has the most significant impact on the study outcomes?",
			OptionA:       "Sample demographics",
			OptionB:       "Geographic location",
			OptionC:       "The analytical framework and data quality",
			OptionD:       "Publication timeline",
			CorrectAnswer: "C",
			Explanation:   "The paper emphasizes that the analytical framework design and the quality of input data are the primary determinants of study outcomes, as demonstrated through sensitivity analysis.",
		},
		{
			Question:      "What is the main limitation acknowledged by the authors?",
			OptionA:       "Insufficient funding for the research",
			OptionB:       "The cross-sectional nature limiting causal inferences",
			OptionC:       "Lack of peer review",
			OptionD:       "Outdated data sources",
			CorrectAnswer: "B",
			Explanation:   "The authors explicitly acknowledge that the cross-sectional design of their study prevents them from making causal claims about the relationships between the variables under investigation.",
		},
		{
			Question:      "How does this paper contribute to the existing body of knowledge?",
			OptionA:       "It replicates previous studies without new findings",
			OptionB:       "It proposes a novel analytical framework with improved performance",
			OptionC:       "It focuses solely on theoretical aspects without empirical validation",
			OptionD:       "It provides only a literature review",
			CorrectAnswer: "B",
			Explanation:   "The paper's primary contribution is the development and validation of a novel analytical framework that demonstrates measurably improved performance compared to existing baseline approaches.",
		},
		{
			Question:      "What type of validation was used to confirm the robustness of the results?",
			OptionA:       "Peer review only",
			OptionB:       "Single dataset testing",
			OptionC:       "Cross-validation across multiple datasets",
			OptionD:       "No validation was performed",
			CorrectAnswer: "C",
			Explanation:   "The authors used cross-validation experiments across multiple datasets to confirm that their results are robust and generalizable beyond the primary dataset.",
		},
		{
			Question:      "Which future research direction is suggested by the authors?",
			OptionA:       "Abandoning the current approach entirely",
			OptionB:       "Conducting longitudinal studies to establish causal relationships",
			OptionC:       "Reducing the scope of the research",
			OptionD:       "Focusing only on theoretical analysis",
			CorrectAnswer: "B",
			Explanation:   "The authors recommend longitudinal studies as a key future direction, which would allow researchers to track how identified patterns evolve over time and enable causal inferences.",
		},
		{
			Question:      "What statistical significance level was used in the study?",
			OptionA:       "p < 0.10",
			OptionB:       "p < 0.01",
			OptionC:       "p < 0.05",
			OptionD:       "No significance testing was performed",
			CorrectAnswer: "C",
			Explanation:   "The study used a standard significance level of p < 0.05 for statistical hypothesis testing, which is the conventional threshold in social science and applied research.",
		},
		{
			Question:      "What approach was used for data collection in this study?",
			OptionA:       "Single-source data collection",
			OptionB:       "Systematic multi-source collection with triangulation",
			OptionC:       "Web scraping only",
			OptionD:       "Government databases exclusively",
			CorrectAnswer: "B",
			Explanation:   "The study employed systematic data collection from multiple sources with triangulation to ensure validity and reduce bias in the findings.",
		},
		{
			Question:      "What is the significance of the cross-validation results?",
			OptionA:       "They show the method only works on one dataset",
			OptionB:       "They confirm the generalizability and robustness of the findings",
			OptionC:       "They indicate the need for more data",
			OptionD:       "They suggest the methodology is flawed",
			CorrectAnswer: "B",
			Explanation:   "The cross-validation results are significant because they confirm that the proposed approach produces consistent results across different datasets, demonstrating both robustness and generalizability.",
		},
		{
			Question:      "How does the paper address potential bias in the research?",
			OptionA:       "It ignores the issue of bias entirely",
			OptionB:       "It uses only automated data collection",
			OptionC:       "It acknowledges bias limitations and uses multiple validation strategies",
			OptionD:       "It relies solely on peer review for bias mitigation",
			CorrectAnswer: "C",
			Explanation:   "The paper transparently acknowledges potential sources of bias (such as recall bias and social desirability effects) and addresses them through methodological choices including data triangulation and cross-validation.",
		},
	}

	// Return the requested number of questions
	result := make([]models.QuizQuestionAI, 0, numQuestions)
	for i := 0; i < numQuestions && i < len(questionPool); i++ {
		result = append(result, questionPool[i])
	}

	return result, nil
}

// ChatWithContext generates a contextual response to a question using paper chunks.
// The mock service provides intelligent responses based on keyword detection.
func (m *MockAIService) ChatWithContext(question string, context string, chatHistory []models.ChatMessage) (string, error) {
	// Simulate processing time
	time.Sleep(200 * time.Millisecond)

	questionLower := strings.ToLower(question)

	// Check for common question patterns and provide relevant responses
	switch {
	case containsAny(questionLower, "summary", "summarize", "overview", "about"):
		return "Based on the paper's content, this research presents a comprehensive study that examines " +
			"key aspects of the topic through rigorous methodology. The authors have conducted " +
			"systematic analysis using both quantitative and qualitative approaches. The main findings " +
			"indicate significant patterns and trends that contribute to our understanding of the field. " +
			"The paper provides evidence-based conclusions supported by robust statistical analysis " +
			"and cross-validated results.", nil

	case containsAny(questionLower, "method", "methodology", "approach", "how did"):
		return "The paper employs a mixed-methods research approach. The quantitative component " +
			"involves statistical analysis including regression modeling, hypothesis testing, and " +
			"cross-validation techniques. The qualitative component uses thematic analysis of " +
			"expert interviews and case studies. Data was collected from multiple sources to ensure " +
			"triangulation and validity. The authors follow established research protocols with " +
			"appropriate ethical guidelines and data quality measures throughout the research process.", nil

	case containsAny(questionLower, "finding", "result", "discover", "conclude", "conclusion"):
		return "The key findings of this paper include: (1) Significant patterns and trends were " +
			"identified that had not been previously documented in the existing literature. " +
			"(2) The proposed analytical framework demonstrates improved performance over baseline " +
			"approaches. (3) Strong correlations were found between the primary variables under " +
			"investigation with statistical significance (p < 0.05). (4) The results have " +
			"practical implications for both industry practitioners and policymakers. " +
			"(5) Cross-validation confirmed the robustness of the approach across multiple datasets.", nil

	case containsAny(questionLower, "limitation", "weakness", "drawback", "problem"):
		return "The authors acknowledge several limitations in their study: First, the sample size, " +
			"while adequate for the methodology, may limit generalizability to broader populations. " +
			"Second, the cross-sectional design prevents causal inferences about relationships between " +
			"variables. Third, self-reported data may be subject to recall bias and social desirability " +
			"effects. Finally, the specific geographic and temporal context may affect external validity. " +
			"These limitations are transparently discussed and provide important context for interpreting the results.", nil

	case containsAny(questionLower, "future", "recommend", "suggestion", "next step"):
		return "The paper suggests several promising directions for future research: (1) Longitudinal " +
			"studies should be conducted to examine how identified patterns evolve over time. " +
			"(2) Replication across different geographic regions and cultural contexts would strengthen " +
			"external validity. (3) Integration of emerging technologies like machine learning and " +
			"big data analytics could enhance the analytical framework. (4) Collaborative " +
			"interdisciplinary research would provide a more holistic understanding of the phenomena. " +
			"These directions represent natural extensions of the current work.", nil

	case containsAny(questionLower, "data", "dataset", "sample"):
		return "The research utilizes data collected through a systematic multi-source approach. " +
			"The dataset includes both quantitative measurements and qualitative observations gathered " +
			"following established protocols. Data quality was ensured through triangulation from " +
			"multiple independent sources. The authors describe their data collection process in detail, " +
			"including sampling strategy, data cleaning procedures, and validation steps to ensure " +
			"reliability and reproducibility of the findings.", nil

	case containsAny(questionLower, "compare", "versus", "difference", "better", "worse"):
		return "The paper provides comparative analysis between the proposed approach and existing " +
			"baseline methods. The results demonstrate that the novel framework achieves measurably " +
			"improved performance across key evaluation metrics. Specifically, the authors report " +
			"improvements in accuracy, robustness, and generalizability when compared to traditional " +
			"approaches. The cross-validation experiments further confirm that these improvements " +
			"are consistent across different datasets and experimental conditions.", nil

	default:
		// Generate a response that references the provided context
		contextPreview := context
		if len(contextPreview) > 200 {
			contextPreview = contextPreview[:200] + "..."
		}

		return fmt.Sprintf(
			"Based on the relevant sections of the paper, I can provide the following insights regarding your question:\n\n"+
				"The paper addresses aspects related to your query through its comprehensive analysis. "+
				"The relevant context from the paper mentions: \"%s\"\n\n"+
				"The authors' analysis suggests that this topic is interconnected with the broader themes "+
				"explored in the research. For a more detailed understanding, I would recommend reviewing "+
				"the methodology and results sections of the paper, which provide quantitative evidence "+
				"supporting the conclusions drawn.",
			contextPreview,
		), nil
	}
}

// containsAny checks if the text contains any of the given keywords.
func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
