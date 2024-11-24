# 🛸 Go Script for Collecting Star Dust in Tiny Universe 🌌

This project provides a 🚀 **bot** that collects star dust from a service using an API and sends updates to **Telegram**. You can configure it to run at random intervals and track your progress over time.

[Русский](README_RU.md) | **English**

---

## 📋 Requirements

- **🛠️ Go version 1.20 or higher**
- **📦 Packages**: Required dependencies are listed in `go.mod`.

---

## 🚀 Getting Started

### 🔍 1. Obtaining Your Session ID and Galaxy ID

#### How to find `GALAXY_ID`:
1. 🖥️ Open your browser and go to [Telegram Web](https://web.telegram.org).
2. 🛠️ Press `F12` to open the Developer Tools.
3. Go to the **"Network"** tab.
4. Perform an action to create stars (it doesn’t matter if you can buy them; you just need to send the request).
5. Find the `POST` request named `create` in the network logs.
6. In the **"Payload"** section of this request, locate the `galaxy_id` field. This is your `GALAXY_ID`.

---

### 🔎 2. Finding Your Telegram ID

To find your `TELEGRAM_ID`, use a Telegram bot such as [@getMyID_tgbot](https://t.me/getMyID_tgbot).

---

### ⚙️ 3. Configuring the Configuration File

Create a `config.json` file with your configuration:

```dotenv
# Telegram Bot Configuration
BOT_TOKEN="<bot-token>"
SESSION_ID=<sessionId1>,<sessionId2>
TELEGRAM_ID=<tg-user-id>

SEND_TO_TELEGRAM=True

TOLERANCE_FROM=600
TOLERANCE_TO=1500
STATS_PER_REQUEST=10

STARS_AUTO_BUY=True
STARS_AUTO_BUY_COUNT=100
GALAXY_ID=<galaxyId1>,<galaxyId2>
```

Replace the placeholder values `<bot-token>`, `<sessionId>`, `<galaxyId>`, and `<tg-user-id>` with your actual data.

---

### 🏃‍♂️ 4. Running the Bot

#### Install dependencies:
```bash
go mod tidy
```

#### Run the bot:
```bash
go run main.go
```

🎉 The bot will start collecting star dust at random intervals and send updates to your Telegram (if configured).

---

## 🛑 Disclaimer

⚠️ By using this bot, you agree that all responsibility for its use lies solely with you.  
The author is not responsible for any damage, consequences, or actions caused by the use of this software. Ensure you understand the potential risks, and use it at your own risk. 🕊️

---