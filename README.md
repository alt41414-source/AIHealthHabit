# AI Habit Tracker

## Steps to run on your machine

1. `git clone https://github.com/alt41414-source/AIHealthHabit.git`
2. Go to the backend directory, find the `.env.example`, and rename it to `.env`
3. Get API key from https://console.groq.com/keys for Groq and put it in the `GROQ_API_KEY` environment variable
4. Generate a really secure secret key and put it in the `JWT_SECRET` environment variable
5. Go to the `frontend` directory and run `npm i`, then run `npm run build`
6. Start the server from the `./server` binary in `backend`, if it doesn't work run `go build -o server` in the `backend` directory
