package fire

import "encoding/json"

var payload = json.RawMessage(pp)

const pp = `
{
  "action": "opened",
  "number": 42,
  "pull_request": {
    "url": "https://api.github.com/repos/username/random-repo/pulls/42",
    "id": 123456789,
    "node_id": "MDExOlB1bGxSZXF1ZXN0MTIzNDU2Nzg5",
    "html_url": "https://github.com/username/random-repo/pull/42",
    "diff_url": "https://github.com/username/random-repo/pull/42.diff",
    "patch_url": "https://github.com/username/random-repo/pull/42.patch",
    "issue_url": "https://api.github.com/repos/username/random-repo/issues/42",
    "number": 42,
    "state": "open",
    "locked": false,
    "title": "Add new feature",
    "user": {
      "login": "contributor_username",
      "id": 999888777,
      "avatar_url": "https://avatars.githubusercontent.com/u/999888777?v=4",
      "html_url": "https://github.com/contributor_username"
    },
    "body": "This pull request adds a fantastic new feature to the project.",
    "created_at": "2023-08-09T12:34:56Z",
    "updated_at": "2023-08-09T12:34:56Z",
    "closed_at": null,
    "merged_at": null,
    "merge_commit_sha": null,
    "assignee": null,
    "assignees": [],
    "requested_reviewers": [],
    "requested_teams": [],
    "labels": [],
    "milestone": null,
    "commits_url": "https://api.github.com/repos/username/random-repo/pulls/42/commits",
    "review_comments_url": "https://api.github.com/repos/username/random-repo/pulls/42/comments",
    "review_comment_url": "https://api.github.com/repos/username/random-repo/pulls/comments{/number}",
    "comments_url": "https://api.github.com/repos/username/random-repo/issues/42/comments",
    "statuses_url": "https://api.github.com/repos/username/random-repo/statuses/{sha}",
    "head": {
      "label": "username:feature-branch",
      "ref": "feature-branch",
      "sha": "abcdef1234567890abcdef1234567890abcdef12",
      "user": {
        "login": "username",
        "id": 123456789
      },
      "repo": {
        "id": 987654321,
        "name": "random-repo",
        "full_name": "username/random-repo",
        "private": false,
        "owner": {
          "login": "username",
          "id": 123456789
        },
        "html_url": "https://github.com/username/random-repo"
      }
    },
    "base": {
      "label": "username:main",
      "ref": "main",
      "sha": "b9e47a3c2f0e8d7b6a8f4e2c1d3b5e6f7a8b9c0d",
      "user": {
        "login": "username",
        "id": 123456789
      },
      "repo": {
        "id": 987654321,
        "name": "random-repo",
        "full_name": "username/random-repo",
        "private": false,
        "owner": {
          "login": "username",
          "id": 123456789
        },
        "html_url": "https://github.com/username/random-repo"
      }
    }
  },
  "repository": {
    "id": 987654321,
    "node_id": "MDEwOlJlcG9zaXRvcnkxMjM0NTY3ODk=",
    "name": "random-repo",
    "full_name": "username/random-repo",
    "private": false,
    "owner": {
      "login": "username",
      "id": 123456789,
      "avatar_url": "https://avatars.githubusercontent.com/u/123456789?v=4",
      "html_url": "https://github.com/username"
    },
    "html_url": "https://github.com/username/random-repo",
    "description": "This is a randomly generated GitHub repository.",
    "fork": false,
    "url": "https://api.github.com/repos/username/random-repo",
    "forks_url": "https://api.github.com/repos/username/random-repo/forks",
    "keys_url": "https://api.github.com/repos/username/random-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/username/random-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/username/random-repo/teams",
    "hooks_url": "https://api.github.com/repos/username/random-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/username/random-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/username/random-repo/events",
    "assignees_url": "https://api.github.com/repos/username/random-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/username/random-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/username/random-repo/tags",
    "blobs_url": "https://api.github.com/repos/username/random-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/username/random-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/username/random-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/username/random-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/username/random-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/username/random-repo/languages",
    "stargazers_url": "https://api.github.com/repos/username/random-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/username/random-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/username/random-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/username/random-repo/subscription",
    "commits_url": "https://api.github.com/repos/username/random-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/username/random-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/username/random-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/username/random-repo/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/username/random-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/username/random-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/username/random-repo/merges",
    "archive_url": "https://api.github.com/repos/username/random-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/username/random-repo/downloads",
    "issues_url": "https://api.github.com/repos/username/random-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/username/random-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/username/random-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/username/random-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/username/random-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/username/random-repo/releases{/id}",
    "deployments_url": "https://api.github.com/repos/username/random-repo/deployments",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-08-09T12:34:56Z",
    "pushed_at": "2023-08-09T12:34:56Z",
    "git_url": "git://github.com/username/random-repo.git",
    "ssh_url": "git@github.com:username/random-repo.git",
    "clone_url": "https://github.com/username/random-repo.git",
    "svn_url": "https://github.com/username/random-repo",
    "homepage": "https://github.com/username/random-repo",
    "size": 1024,
    "stargazers_count": 42,
    "watchers_count": 42,
    "language": "Python",
    "has_issues": true,
    "has_projects": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": true,
    "forks_count": 10,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 5,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit"
    },
    "forks": 10,
    "open_issues": 5,
    "watchers": 42,
    "default_branch": "main"
  },
  "organization": {
    "login": "org_username",
    "id": 777666555,
    "url": "https://api.github.com/orgs/org_username",
    "repos_url": "https://api.github.com/orgs/org_username/repos",
    "events_url": "https://api.github.com/orgs/org_username/events",
    "hooks_url": "https://api.github.com/orgs/org_username/hooks",
    "issues_url": "https://api.github.com/orgs/org_username/issues",
    "members_url": "https://api.github.com/orgs/org_username/members{/member}",
    "public_members_url": "https://api.github.com/orgs/org_username/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/777666555?v=4",
    "description": "An imaginary GitHub organization"
  },
  "sender": {
    "login": "sender_username",
    "id": 111222333,
    "avatar_url": "https://avatars.githubusercontent.com/u/111222333?v=4",
    "html_url": "https://github.com/sender_username"
  }
}

`
