# 📡 GoSniff

![250317_12h44m11s_screenshot](https://github.com/user-attachments/assets/77302853-b5ca-4228-a654-baa245dd786d)

An RSS feed aggregator built with Go that helps you stay updated with your favorite content! 🚀

## ✨ Features

- 🔐 User Authentication via API Keys
- 📂 Feed Management (Create, List, Follow, Unfollow)
- 🤝 Multi-user Support
- 🔄 Automatic Feed Aggregation
- 🕒 Real-time Updates
- 🎯 REST API Interface

## 🛠️ Technologies

- Go 1.x
- PostgreSQL
- Chi Router
- SQLC
- Goose (migrations)

## 🏗️ Project Architecture
📦 GoSniff
 ┣ 📂 internal
 ┃ ┗ 📂 database         # Database models and queries
 ┣ 📂 sql
 ┃ ┣ 📂 schema           # Database migrations
 ┃ ┗ 📂 queries          # SQLC queries
 ┗ 📜 main.go            # Application entrypoint

## 🚀 Getting Started

1. Clone the repository
```bash
git clone https://github.com/WST-T/GoSniff.git
```

2. Set up your environment variables
```bash
.env
# Edit .env with your configuration
It should be something like:
PORT=?
DB_URL=?
```

3. Run database migrations
```bash
goose postgres "postgres://user:password@localhost:5432/gosniffdb" up
goose postgres "postgres://user:password@localhost:5432/gosniffdb" down
```

4. Start the server
```bash
go run main.go ||
go build && ./GoSniff
```

## 🔄 Feed Aggregator
The feed aggregator worker runs in the background and:

- ⏰ Periodically checks RSS feeds for updates
- 📥 Fetches new content
- 💾 Stores new posts in the database
- 🔔 Keeps content fresh and up-to-date

## 🤝 Contributing
Feel free to open issues and pull requests!

## 📝 License
MIT License - feel free to use this project however you'd like!

## 🙏 Acknowledgments
- Chi router for the amazing HTTP routing
- SQLC for type-safe SQL
- Goose for easy database migrations
- And all other open source contributors!

Made with ❤️ using Go