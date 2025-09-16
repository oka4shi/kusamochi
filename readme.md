# Kusamochi

Kusamochi (草餅, rice cakes mixed with mugwort in Japanese) is a tool designed to compile GitHub contributions from the previous week and send out a ranking via Discord Webhook on a weekly basis.

## Usage

### Run Locally

1. Generate a new Personal Access Token (classic) on GitHub. The `user` scope is required.
2. Duplicate `.templ.env.sh`, and complete the environment variables.
3. Create a file named `users.json` and list the the GitHub usernames of the person whose contributions you want to track in the following format:

```json
["username1", "username2"]
```

4. Build the application using `make build`.
5. Run `source .env.sh` and then `./kusamochi` in the shell.

### Using Docker

1. Generate a new Personal Access Token (classic) on GitHub. The `user` scope is required.
2. Duplicate `.templ.env`, and complete the environment variables.
3. Create a file named `users.json` and list the the GitHub usernames of the person whose contributions you want to track in the following format:
4. Build a docker image: `docker build -t kusamochi .`
5. Run the docker image: `docker run --env-file .env --name kusamochi --volume <abosolute-path-of-your-users.json-file>:/app/users.json --rm -it kusamochi`

## License
The original idea of this application is designed by @s7tya. I'm grateful for his contribution.

This application is licensed under the MIT License. 

