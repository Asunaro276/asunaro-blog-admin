resource "aws_ecr_repository" "default" {
  name                 = var.image_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "null_resource" "build_and_push_image" {
  depends_on = [aws_ecr_repository.default]
  provisioner "local-exec" {
    // ../../cms_api/${var.language}のdeploy.shを実行 引数としてimage_name, language, repository_url, region_nameを渡す
    command = "bash ../../cms_api/${var.language}/deploy.sh ${data.external.git.result.sha} ${var.language} ${aws_ecr_repository.default.repository_url} ${data.aws_region.current_region.name}"
  }

  triggers = {
    // codedir_local_path配下のファイルを結合してsha256を計算したものを使って差分検知を行う
    code_diff = join("", [
      for file in fileset("${local.codedir_local_path}", "{**}") : filebase64sha256("${local.codedir_local_path}/${file}")
    ])
  }
}
