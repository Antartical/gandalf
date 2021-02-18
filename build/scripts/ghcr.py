"""
Configures ghcr credentials in your local machine.
"""


import os


HOME = os.getenv('HOME')


def write_files(username: str, token: str):
    if not os.path.exists(f'{HOME}/.credentials'):
        os.mkdir(f'{HOME}/.credentials')
    with open(f'{HOME}/.credentials/ghcr.name', 'w') as ghcr_name:
        ghcr_name.write(username)
    with open(f'{HOME}/.credentials/ghcr.token', 'w') as ghcr_token:
        ghcr_token.write(token)


if __name__ == "__main__":
    username = input("Enter your github username: ")
    token = input("Enter your github token: ")
    write_files(username, token)
