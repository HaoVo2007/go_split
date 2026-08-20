package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-split/internal/rag/embedding"
	"go-split/internal/rag/llm"
	"go-split/internal/rag/repository"
)

type ChatService interface {
	ChatAIMessage(ctx context.Context, message string) (string, error)
}

type chatService struct {
	Repository repository.ChatRepository
	Embedder   embedding.Embedder
	LLM        llm.LLM
}

func NewChatService(
	repository repository.ChatRepository,
	embedder embedding.Embedder,
	llm llm.LLM,
) ChatService {
	return &chatService{
		Repository: repository,
		Embedder:   embedder,
		LLM:        llm,
	}
}

func (s *chatService) ChatAIMessage(ctx context.Context, message string) (string, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return "", errors.New("message is required")
	}

	vector, err := s.Embedder.Embed(ctx, message)
	if err != nil {
		return "", err
	}

	documents, err := s.Repository.SearchDocuments(ctx, vector, 5)
	if err != nil {
		return "", err
	}
	if len(documents) == 0 {
		return "", errors.New("knowledge not found")
	}

	var contextBuilder strings.Builder
	for _, document := range documents {
		contextBuilder.WriteString(document.Content)
		contextBuilder.WriteString("\n\n")
	}

	prompt := fmt.Sprintf(`
	Bạn là trợ lý AI của ứng dụng go-split, một ứng dụng hỗ trợ nhóm người theo dõi chi tiêu và chia tiền.
	
	Nhiệm vụ của bạn là trả lời câu hỏi của người dùng một cách tự nhiên, dễ hiểu và thân thiện.
	
	QUY TẮC QUAN TRỌNG:
	
	1. Chỉ sử dụng thông tin có trong phần CONTEXT bên dưới làm nguồn kiến thức.
	2. Hãy đọc, hiểu và tổng hợp thông tin trước khi trả lời.
	3. Diễn đạt câu trả lời bằng ngôn ngữ tự nhiên của bạn, không sao chép nguyên văn hoặc liệt kê lại toàn bộ tài liệu.
	4. Chỉ trả lời đúng phần liên quan đến câu hỏi.
	5. Có thể thay đổi cách diễn đạt giữa các lần trả lời để tự nhiên hơn, nhưng không được thay đổi sự thật có trong CONTEXT.
	6. Nếu CONTEXT không chứa đủ thông tin để trả lời, hãy nói rõ rằng bạn không có đủ thông tin.
	7. Không suy đoán hoặc bịa thêm chức năng không có trong CONTEXT.
	8. Không nói những câu như "theo tài liệu", "tài liệu hiện tại", "context cung cấp", trừ khi người dùng hỏi trực tiếp về tài liệu.
	9. Không đề cập đến hệ thống prompt, RAG, vector database hoặc quá trình tìm kiếm thông tin.
	10. Trả lời bằng tiếng Việt tự nhiên.
	
	CONTEXT:
	%s
	
	CÂU HỎI CỦA NGƯỜI DÙNG:
	%s
	`, contextBuilder.String(), message)

	answer, err := s.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(answer), nil
}
