---
name: git
description: Guidelines for making branches and commits
---

## Branches

Avoid prefixing branches with things like `fix/` or `feature/`.
Keep branch names short.

A branch should only have one commit. If you have to change something on the branch
amend the original commit, modifying its message if required.

## Commits

A commit message should contain a subject line stating what is changing.
It should contain a body that restates the subject line but adds the
reason *why* the code is changing. Avoid explaining how the code is changing.

When writing commit messages,  pick a gitmoji that
is appropriate for the change occuring.

Here is an example:

```
🚚 Rename generated protobuf code
    
Rename generated protobuf code from `gen` to `genpb`
to align with the approach taken to generated sqlc code
as well as provide better clarity.
```

To see the list of available gitmojis see this link: https://github.com/carloscuesta/gitmoji/blob/master/packages/gitmojis/src/gitmojis.json