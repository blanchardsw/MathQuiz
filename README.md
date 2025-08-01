# Mental Math Trainer

A full-stack application for improving mental math skills with timed exercises.

## Architecture

- **Frontend**: React + TypeScript with Tailwind CSS
- **Backend**: Go with Gorilla Mux router
- **Communication**: REST API

## Project Structure

```
mental-math-trainer/
├── backend/
│   ├── main.go              # Server entry point
│   ├── handlers/
│   │   └── quiz.go          # API endpoints
│   ├── models/
│   │   └── question.go      # Data structures
│   └── utils/
│       └── generator.go     # Question generation logic
├── frontend/
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── components/
│   │   │   ├── Quiz.tsx     # Main quiz interface
│   │   │   ├── Timer.tsx    # Timer component
│   │   │   └── Score.tsx    # Score display
│   │   ├── pages/
│   │   │   └── Home.tsx     # Main page
│   │   ├── services/
│   │   │   └── api.ts       # API communication
│   │   └── App.tsx
│   └── package.json
├── README.md
└── go.mod
```

## API Endpoints

- `GET /api/quiz` - Get a new math question
- `POST /api/answer` - Submit an answer
- `GET /api/score` - Get current session score

## Getting Started

### Backend Setup

1. Navigate to the project root:
   ```bash
   cd "c:\Users\blanc\CascadeProjects\Math Trainer"
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the backend server:
   ```bash
   go run backend/main.go
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The frontend will start on `http://localhost:3000`

## Features

- **Random Question Generation**: Addition, subtraction, multiplication, and division
- **Timer**: 30-second countdown per question
- **Score Tracking**: Real-time accuracy and progress tracking
- **Responsive UI**: Modern design with Tailwind CSS
- **Visual Feedback**: Color-coded timer and accuracy indicators

## Technologies Used

- **Frontend**: React 18, TypeScript, Tailwind CSS, Axios
- **Backend**: Go, Gorilla Mux, CORS middleware
- **Development**: Hot reload, TypeScript compilation

## Next Steps

- Add difficulty levels
- Implement user authentication
- Add persistent score storage
- Create leaderboards
- Add more question types
