from github import Github, Auth, GithubException
from github.Commit import Commit

import os
import sys
import itertools
from datetime import date, datetime, timezone, timedelta


def is_on_day(commit: Commit, day: date):
    # Make sure it's the correct type of object.
    commit_date: date = commit.commit.author.date.date()
    return commit_date == day


def main():
    token = os.getenv("GITHUB_PAT")
    assert token

    auth = Auth.Token(token)

    day_ago = datetime.now(timezone.utc) - timedelta(days=1)
    today: date = datetime.today().date()
    print(f"Date: {today}")
    print(f"Timezone: {day_ago.tzinfo}")

    total_priv = 0
    total_pub = 0
    total_today_commits = 0

    with Github(auth=auth) as g:
        for repo in g.get_user().get_repos():
            err: None | str = None
            today_commits: list[Commit] = []

            try:
                # This should return commits only from the default branch.
                commits = repo.get_commits(since=day_ago)
                latest = itertools.takewhile(lambda c : is_on_day(c, today), commits)
                today_commits = list(latest)
            except GithubException as e:
                err = e.data["message"]

            if repo.private:
                total_priv += 1
            else:
                total_pub += 1

            if today_commits or err is not None:
                total_today_commits += len(today_commits)
                priv = '(public)' if not repo.private else ''
                print(f"\n{repo.name} {priv}")

                if err:           print(f"  | ERROR: {err}")
                if today_commits:
                    for c in today_commits:
                        print(f"  | {c.commit.message}")

    print(f"\nRepositories checked: {total_priv + total_pub} (public={total_pub}, private={total_priv})")
    print(f"Today commit count: {total_today_commits}")

    if total_today_commits > 0:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
