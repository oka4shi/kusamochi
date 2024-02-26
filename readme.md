# Kusamochi

Kusamochi (草餅, rice cakes mixed with mugwort in Japanese) is a tool designed to compile GitHub contributions from the previous week and send out a ranking via Discord Webhook on a weekly basis.

## Usage

1. Generate a new Personal Access Token (classic) on GitHub. The `user` scope is required.
2. Duplicate `.templ.env.sh`, and complete the environment variables.
3. Create a file named `users.json` and list the the GitHub usernames of the person whose contributions you want to track in the following format:

```json
["username1", "username2"]
```

4. Build the application using `go build`.
5. Run `source .env.sh` and then `./kusamochi` in the shell.

## License
The original idea of this application is designed by @s7tya. I'm grateful for his contribution.

This application is licensed under the MIT License. 

