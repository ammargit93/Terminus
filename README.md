# Terminus — Terminal AI Chat with RAG

Terminus is a terminal-based AI chat application built in Go that integrates large language models (LLMs) with retrieval-augmented generation (RAG). It allows users to query documents, retrieve contextually relevant information via embeddings, and receive AI-generated responses, all from a fully interactive terminal interface.

## 🚀 Features

- **TUI (Terminal User Interface)**
  - Interactive terminal UI using Bubble Tea
  - Scrollable chat viewport
  - Input chatbox with focus management
  - File picker overlay for selecting multiple documents

- **RAG (Retrieval-Augmented Generation)**
  - Reads and stores file content embeddings
  - Computes cosine similarity between user queries and file embeddings
  - Retrieves top-K relevant documents for context-aware responses

- **LLM Integration**
  - Custom Go wrapper for Groq LLM API
  - Sends user queries and retrieved content to generate intelligent responses

- **Vector Store**
  - In-memory storage of embeddings and associated file paths
  - Add, remove, and retrieve embedding pairs easily

- **Modular and Extensible**
  - Embeddings and retrieval logic can be extended to handle chunked files
  - Architecture supports adding agentic behavior for autonomous multi-step reasoning

### To get started:

For Linux and MacOS
```bash
export COHERE_API_KEY=<your-cohere-api-key>
export GROQ_API_KEY=<your-groq-api-key>
```
For windows (powershell) 
```bash
setx COHERE_API_KEY <your-cohere-api-key>
setx GROQ_API_KEY <your-groq-api-key>
```
restart your IDE after configuring

**Clone the repository**
```bash
git clone https://github.com/ammargit93/Terminus.git
cd Terminus
```

Run the program
```bash
go run .
```
## ⚡ Usage

- Type messages in the chatbox and press Enter to send.

- Press Esc to open the file picker and select documents for context.

- Use arrow keys or PgUp/PgDown to scroll through chat history.

- Press Tab to toggle chatbox focus.

- Press Ctrl+C to quit the application.