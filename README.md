<h1 align="center">
  Accountability Discord Bot
  <br>
</h1>

<p align="center">
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg"></a>
  <a href="https://twitter.com/adamjsturge"><img src="https://img.shields.io/twitter/follow/adamjsturge.svg?logo=twitter"></a>
</p>

# Setup

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Discord Bot Configuration
DISCORD_TOKEN=your_discord_bot_token_here

# PostgreSQL Database Configuration
POSTGRES_HOST=pg_db
POSTGRES_PORT=5432
POSTGRES_USER=accountability_user
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=accountability_optcg

# Database Connection Retry Configuration (optional)
DB_MAX_RETRIES=5
DB_RETRY_DELAY_SECONDS=10
```

# To Add before release

- [x] Database Postgres
- [x] Tables for users
- [ ] Accountability Basics: Set Timezone
- [ ] Accountability Basics: Upload your matches results (either 1 at time or multiple -> each entry gets its own row)
- [ ] Auth Basics: Tags for permissions
- [ ] Accountability Basics: Categories for practice (Locals, Ranked, etc)
- [ ] Accountability Basics: Streak tracking (consecutive days practicing)

# To Add after release

- [ ] Accountability Advanced: Daily/Weekly Reminders
- [ ] Accountability Advanced: Daily/Weekly Summary (Opt-InS)
- [ ] Accountability Advanced: Goal setting & tracking (e.g., "Play 5 games this week")
- [ ] Accountability Advanced: Public leaderboards (most active players)
- [ ] Accountability Advanced: Teams within a server
- [ ] Accountability Advanced: Team Goals

# Stretch Goal

- [ ] Say what tournaments you plan to go to and get reminders when they are coming up and when to practice for them
- [ ] Last thing: See if you can hook into ranked matches or kaizoku to auto upload your matches results
