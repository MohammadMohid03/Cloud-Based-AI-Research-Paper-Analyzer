package database

import (
	"log"
	"time"

	"github.com/research-paper-analyzer/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RunMigrations performs GORM auto-migration for all models.
// This creates or updates tables to match the current model definitions.
// In production, you might want to use a proper migration tool like golang-migrate.
func RunMigrations() error {
	log.Println("🔄 Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Paper{},
		&models.PaperAnalysis{},
		&models.PaperChunk{},
		&models.ChatHistory{},
		&models.QuizQuestion{},
	)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// SeedData populates the database with demo data for development and testing.
// It creates a demo user and a sample paper with pre-populated analysis results.
// This function is idempotent - it won't create duplicates if run multiple times.
func SeedData() error {
	log.Println("🌱 Seeding database with demo data...")

	// --- Seed Demo User ---
	var existingUser models.User
	result := DB.Where("email = ?", "demo@example.com").First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		// Hash the demo password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("demo123"), 10)
		if err != nil {
			return err
		}

		demoUser := models.User{
			Name:         "Demo User",
			Email:        "demo@example.com",
			PasswordHash: string(hashedPassword),
		}

		if err := DB.Create(&demoUser).Error; err != nil {
			return err
		}

		log.Println("  ✅ Demo user created (demo@example.com / demo123)")

		// --- Seed Demo Paper ---
		demoPaper := models.Paper{
			UserID:  demoUser.ID,
			Title:   "Deep Learning Approaches for Natural Language Processing: A Comprehensive Survey",
			FileURL: "demo-paper.pdf",
			Status:  models.StatusCompleted,
			RawText: getSamplePaperText(),
		}

		if err := DB.Create(&demoPaper).Error; err != nil {
			return err
		}

		log.Println("  ✅ Demo paper created")

		// --- Seed Paper Analysis ---
		demoAnalysis := models.PaperAnalysis{
			PaperID: demoPaper.ID,
			Summary: "This comprehensive survey examines the evolution and current state of deep learning approaches in Natural Language Processing (NLP). The paper traces the progression from traditional statistical methods to modern transformer-based architectures, highlighting key breakthroughs such as BERT, GPT, and their variants. The authors provide a systematic comparison of various model architectures across multiple NLP tasks including text classification, named entity recognition, machine translation, and question answering. The survey identifies critical challenges in the field including computational costs, data requirements, model interpretability, and bias in language models, while also discussing promising future research directions.",
			KeyFindings: `• Transformer-based models have achieved state-of-the-art results across virtually all NLP benchmarks, surpassing previous RNN and CNN-based approaches.
• Pre-training on large unlabeled corpora followed by task-specific fine-tuning has become the dominant paradigm in NLP.
• Models like BERT and GPT-3 demonstrate emergent abilities when scaled to billions of parameters.
• Transfer learning significantly reduces the amount of labeled data needed for downstream tasks.
• Attention mechanisms are the key innovation enabling transformers to capture long-range dependencies in text.`,
			Methodology: "The authors conducted a systematic literature review covering 150+ papers published between 2017 and 2024. They categorized approaches by architecture type (RNN, CNN, Transformer, Hybrid), pre-training strategy (masked language modeling, causal language modeling, sequence-to-sequence), and application domain. Quantitative comparisons were made using standard benchmarks including GLUE, SuperGLUE, SQuAD, and WMT translation tasks. The review follows PRISMA guidelines for systematic reviews.",
			Limitations: `• The survey primarily focuses on English-language NLP, with limited coverage of multilingual and low-resource language scenarios.
• Computational cost comparisons are limited due to inconsistent reporting across papers.
• The rapid pace of advancement means some very recent models may not be included.
• The paper does not include original experimental results, relying instead on reported benchmarks.`,
			FutureScope: `• Development of more efficient architectures that maintain performance while reducing computational requirements.
• Greater focus on multilingual and cross-lingual transfer learning for low-resource languages.
• Improving model interpretability and explainability for safety-critical applications.
• Addressing bias and fairness in large language models through improved training data curation and debiasing techniques.`,
			Keywords: "deep learning, natural language processing, transformers, BERT, GPT, transfer learning, attention mechanism, pre-training",
		}

		if err := DB.Create(&demoAnalysis).Error; err != nil {
			return err
		}

		log.Println("  ✅ Demo analysis created")

		// --- Seed Paper Chunks ---
		chunks := getSampleChunks(demoPaper.ID)
		for _, chunk := range chunks {
			if err := DB.Create(&chunk).Error; err != nil {
				return err
			}
		}

		log.Println("  ✅ Demo paper chunks created")

		// --- Seed Quiz Questions ---
		quizQuestions := getSampleQuizQuestions(demoPaper.ID)
		for _, q := range quizQuestions {
			if err := DB.Create(&q).Error; err != nil {
				return err
			}
		}

		log.Println("  ✅ Demo quiz questions created")

		// --- Seed Chat History ---
		chatHistory := []models.ChatHistory{
			{
				PaperID:   demoPaper.ID,
				UserID:    demoUser.ID,
				Question:  "What is the main contribution of this paper?",
				Answer:    "The main contribution of this paper is providing a comprehensive and systematic survey of deep learning approaches for Natural Language Processing. It traces the evolution from traditional statistical methods to modern transformer-based architectures like BERT and GPT, offers quantitative comparisons across standard benchmarks, and identifies key challenges and future research directions in the field.",
				CreatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				PaperID:   demoPaper.ID,
				UserID:    demoUser.ID,
				Question:  "What are transformers?",
				Answer:    "Transformers are a type of neural network architecture introduced in the seminal paper 'Attention Is All You Need' by Vaswani et al. (2017). Unlike previous architectures like RNNs and LSTMs that process text sequentially, transformers use a self-attention mechanism to process all words in a sequence simultaneously. This parallel processing enables them to capture long-range dependencies more effectively and train much faster on modern hardware. Transformers form the foundation of modern models like BERT, GPT, and T5.",
				CreatedAt: time.Now().Add(-30 * time.Minute),
			},
		}

		for _, chat := range chatHistory {
			if err := DB.Create(&chat).Error; err != nil {
				return err
			}
		}

		log.Println("  ✅ Demo chat history created")
	} else {
		log.Println("  ℹ️  Demo data already exists, skipping seed")
	}

	log.Println("✅ Database seeding completed")
	return nil
}

// getSamplePaperText returns realistic sample text for the demo paper.
func getSamplePaperText() string {
	return `Deep Learning Approaches for Natural Language Processing: A Comprehensive Survey

Abstract

Natural Language Processing (NLP) has undergone a revolutionary transformation with the advent of deep learning techniques. This comprehensive survey examines the evolution of deep learning approaches in NLP, from early word embeddings to modern transformer-based architectures. We systematically review over 150 papers published between 2017 and 2024, covering key areas including text classification, named entity recognition, machine translation, question answering, and text generation. Our analysis reveals that transformer-based models, particularly pre-trained language models like BERT and GPT, have established new state-of-the-art results across virtually all NLP benchmarks. We identify critical challenges including computational costs, data requirements, model interpretability, and bias, while discussing promising future research directions.

1. Introduction

The field of Natural Language Processing has witnessed unprecedented progress in recent years, driven primarily by advances in deep learning. Traditional NLP approaches relied heavily on hand-crafted features and statistical methods such as n-gram models, Hidden Markov Models, and Conditional Random Fields. While these methods achieved reasonable performance on many tasks, they were limited by their inability to capture complex semantic relationships and long-range dependencies in text.

The introduction of word embeddings, particularly Word2Vec and GloVe, marked the first significant shift toward distributed representations of language. These dense vector representations captured semantic similarities between words and enabled neural networks to process text more effectively. However, these static embeddings could not capture context-dependent word meanings.

The development of recurrent neural networks (RNNs) and their variants, including Long Short-Term Memory (LSTM) networks and Gated Recurrent Units (GRUs), brought the ability to process sequential data and capture temporal dependencies. These architectures became the foundation for many NLP systems, including machine translation, text summarization, and sentiment analysis.

2. The Transformer Revolution

The introduction of the Transformer architecture by Vaswani et al. in their seminal 2017 paper "Attention Is All You Need" fundamentally changed the landscape of NLP. The key innovation was the self-attention mechanism, which allows the model to weigh the importance of different words in a sequence when processing each word. Unlike RNNs, which process sequences step by step, transformers can process all positions simultaneously, enabling much more efficient training on modern parallel hardware.

The multi-head attention mechanism extends this concept by allowing the model to attend to information from different representation subspaces at different positions. This enables the model to capture various types of relationships between words simultaneously. The combination of self-attention with position-wise feed-forward networks, residual connections, and layer normalization creates a powerful architecture for sequence modeling.

3. Pre-trained Language Models

The concept of pre-training large language models on massive text corpora and then fine-tuning them on specific downstream tasks has become the dominant paradigm in NLP. This approach leverages transfer learning to achieve strong performance even with limited task-specific labeled data.

BERT (Bidirectional Encoder Representations from Transformers) introduced masked language modeling as a pre-training objective, where the model learns to predict randomly masked words in a sentence using both left and right context. This bidirectional approach proved highly effective for understanding tasks such as text classification, named entity recognition, and question answering.

GPT (Generative Pre-trained Transformer) and its successors take a different approach, using causal language modeling where the model predicts the next token in a sequence. GPT-3 demonstrated that scaling to 175 billion parameters enables few-shot and zero-shot learning, where the model can perform tasks with minimal or no task-specific training examples.

4. Methodology

We conducted a systematic literature review following PRISMA guidelines. Our search covered major databases including IEEE Xplore, ACM Digital Library, arXiv, and Google Scholar. We identified over 500 potentially relevant papers, which were filtered to 150 core papers based on relevance, citation count, and recency. Each paper was categorized by architecture type, pre-training strategy, and application domain.

5. Limitations and Future Directions

Despite remarkable progress, several challenges remain. The computational cost of training large language models is substantial, often requiring hundreds of GPU-hours and significant energy consumption. Data requirements for pre-training are massive, raising concerns about data quality and bias. Model interpretability remains a significant challenge, particularly for safety-critical applications.

Future research directions include developing more efficient architectures, improving multilingual capabilities, addressing bias and fairness, and enhancing model interpretability. The integration of multimodal learning, combining text with images and other modalities, represents another promising frontier.

6. Conclusion

Deep learning has fundamentally transformed Natural Language Processing, with transformer-based architectures establishing new paradigms for language understanding and generation. The pre-train and fine-tune approach has democratized access to powerful NLP capabilities. However, addressing the remaining challenges of efficiency, bias, and interpretability will be crucial for the responsible deployment of these technologies.`
}

// getSampleChunks returns sample text chunks for the demo paper.
func getSampleChunks(paperID uint) []models.PaperChunk {
	return []models.PaperChunk{
		{
			PaperID:    paperID,
			ChunkIndex: 0,
			ChunkText:  "Deep Learning Approaches for Natural Language Processing: A Comprehensive Survey. Abstract: Natural Language Processing (NLP) has undergone a revolutionary transformation with the advent of deep learning techniques. This comprehensive survey examines the evolution of deep learning approaches in NLP, from early word embeddings to modern transformer-based architectures. We systematically review over 150 papers published between 2017 and 2024, covering key areas including text classification, named entity recognition, machine translation, question answering, and text generation. Our analysis reveals that transformer-based models, particularly pre-trained language models like BERT and GPT, have established new state-of-the-art results across virtually all NLP benchmarks.",
		},
		{
			PaperID:    paperID,
			ChunkIndex: 1,
			ChunkText:  "1. Introduction: The field of Natural Language Processing has witnessed unprecedented progress in recent years, driven primarily by advances in deep learning. Traditional NLP approaches relied heavily on hand-crafted features and statistical methods such as n-gram models, Hidden Markov Models, and Conditional Random Fields. The introduction of word embeddings, particularly Word2Vec and GloVe, marked the first significant shift toward distributed representations of language. The development of recurrent neural networks (RNNs) and their variants, including Long Short-Term Memory (LSTM) networks and Gated Recurrent Units (GRUs), brought the ability to process sequential data and capture temporal dependencies.",
		},
		{
			PaperID:    paperID,
			ChunkIndex: 2,
			ChunkText:  "2. The Transformer Revolution: The introduction of the Transformer architecture by Vaswani et al. in their seminal 2017 paper 'Attention Is All You Need' fundamentally changed the landscape of NLP. The key innovation was the self-attention mechanism, which allows the model to weigh the importance of different words in a sequence when processing each word. Unlike RNNs, which process sequences step by step, transformers can process all positions simultaneously, enabling much more efficient training on modern parallel hardware. The multi-head attention mechanism extends this concept by allowing the model to attend to information from different representation subspaces.",
		},
		{
			PaperID:    paperID,
			ChunkIndex: 3,
			ChunkText:  "3. Pre-trained Language Models: The concept of pre-training large language models on massive text corpora and then fine-tuning them on specific downstream tasks has become the dominant paradigm in NLP. BERT introduced masked language modeling as a pre-training objective, where the model learns to predict randomly masked words using both left and right context. GPT and its successors use causal language modeling where the model predicts the next token in a sequence. GPT-3 demonstrated that scaling to 175 billion parameters enables few-shot and zero-shot learning.",
		},
		{
			PaperID:    paperID,
			ChunkIndex: 4,
			ChunkText:  "4. Methodology: We conducted a systematic literature review following PRISMA guidelines. Our search covered major databases including IEEE Xplore, ACM Digital Library, arXiv, and Google Scholar. We identified over 500 potentially relevant papers, which were filtered to 150 core papers based on relevance, citation count, and recency. 5. Limitations and Future Directions: The computational cost of training large language models is substantial. Data requirements for pre-training are massive, raising concerns about data quality and bias. Future research directions include developing more efficient architectures, improving multilingual capabilities, and addressing bias and fairness.",
		},
	}
}

// getSampleQuizQuestions returns sample quiz questions for the demo paper.
func getSampleQuizQuestions(paperID uint) []models.QuizQuestion {
	return []models.QuizQuestion{
		{
			PaperID:       paperID,
			Question:      "What is the key innovation introduced by the Transformer architecture?",
			OptionA:       "Convolutional layers for text processing",
			OptionB:       "Self-attention mechanism for parallel processing of sequences",
			OptionC:       "Recurrent connections for sequential data",
			OptionD:       "Word2Vec embeddings for word representation",
			CorrectAnswer: "B",
			Explanation:   "The Transformer architecture introduced the self-attention mechanism, which allows the model to process all positions in a sequence simultaneously, unlike RNNs which process sequentially. This enables more efficient training on parallel hardware.",
		},
		{
			PaperID:       paperID,
			Question:      "What pre-training objective does BERT use?",
			OptionA:       "Next sentence prediction only",
			OptionB:       "Causal language modeling (predicting next token)",
			OptionC:       "Masked language modeling (predicting masked words)",
			OptionD:       "Image-text contrastive learning",
			CorrectAnswer: "C",
			Explanation:   "BERT uses masked language modeling as its primary pre-training objective, where random words in a sentence are masked and the model learns to predict them using both left and right context (bidirectional).",
		},
		{
			PaperID:       paperID,
			Question:      "How many parameters does GPT-3 have?",
			OptionA:       "1.5 billion",
			OptionB:       "110 million",
			OptionC:       "175 billion",
			OptionD:       "340 million",
			CorrectAnswer: "C",
			Explanation:   "GPT-3 was scaled to 175 billion parameters, which enabled remarkable few-shot and zero-shot learning capabilities, allowing the model to perform tasks with minimal or no task-specific training examples.",
		},
		{
			PaperID:       paperID,
			Question:      "Which of the following is NOT mentioned as a limitation in the survey?",
			OptionA:       "High computational costs for training",
			OptionB:       "Lack of available training data",
			OptionC:       "Model interpretability challenges",
			OptionD:       "Bias in language models",
			CorrectAnswer: "B",
			Explanation:   "The survey mentions high computational costs, model interpretability challenges, and bias as limitations. However, the issue is not a lack of data but rather concerns about data quality and bias in the massive pre-training datasets.",
		},
		{
			PaperID:       paperID,
			Question:      "What research methodology did the authors follow for their literature review?",
			OptionA:       "Meta-analysis with original experiments",
			OptionB:       "PRISMA guidelines for systematic reviews",
			OptionC:       "Randomized controlled trial",
			OptionD:       "Case study approach",
			CorrectAnswer: "B",
			Explanation:   "The authors conducted a systematic literature review following PRISMA (Preferred Reporting Items for Systematic Reviews and Meta-Analyses) guidelines, reviewing over 150 papers from major academic databases.",
		},
	}
}
