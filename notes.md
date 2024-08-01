# Game

game.ID
game.PeriodCount
game.CurrentUserId

# currentUserId

# Players

_Sorted_

Players len
player.User.ProfileRoot
player.User.Name

--- CardList

templates.IsCurrentUserTurn(ctx)

<!-- /var/repo/stockmarket-new.git/hooks/post-recieve -->

# !/bin/sh
echo "running hook"

whoami

# Print the contents of the working directory

ls -la /var/www/stockmarket-new

# Check out the repository

git --work-tree=/var/www/stockmarket-new --git-dir=/var/repo/stockmarket-new.git checkout -f main

# Print the status after checkout

git --work-tree=/var/www/stockmarket-new --git-dir=/var/repo/stockmarket-new.git status

# Clear the UF cache

rm -rf /var/www/stockmarket-new/app/cache/*

docker-compose -f /var/www/stockmarket-new/backend/docker-compose.yml up -d --build
