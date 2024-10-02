cd ${GIT_REPO_NAME} \
&& git branch -a | grep -v HEAD | perl -ne 'chomp($_); s|^\*?\s*||; if (m|(.+)/(.+)| && not $d{$2}) {print qq(git branch --track $2 $1/$2\n)} else {$d{$_}=1}' | csh -xfs \
&& git fetch --all \
&& git pull --all \
&& git remote set-url origin ${GIT_URL}/${GIT_REPO_NAME}.git \
&& git push -u origin --all \
&& git push -u origin --tags \
&& cd ..
