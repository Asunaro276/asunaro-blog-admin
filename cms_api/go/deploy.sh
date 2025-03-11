# 引数としてimage_name, language, repository_url, region_nameを受け取る
git_sha=$1
language=$2
repository_url=$3
region_name=$4

aws ecr get-login-password --region ${region_name} | docker login --username AWS --password-stdin ${repository_url}

docker build -t ${repository_url}:${git_sha} ../../cms_api/${language}

docker push ${repository_url}:${git_sha}
