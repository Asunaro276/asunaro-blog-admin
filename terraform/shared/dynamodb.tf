# definition.jsonからデータを読み込む
locals {
  definition_json = jsondecode(file("${path.module}/definition.json"))
  data_model      = local.definition_json.DataModel[0]
  table_name      = local.data_model.TableName
  table_data      = local.data_model.TableData

  # パーティションキー/ソートキーの情報
  partition_key = local.data_model.KeyAttributes.PartitionKey
  sort_key      = local.data_model.KeyAttributes.SortKey

  # GSIの情報（複数想定）
  gsis = local.data_model.GlobalSecondaryIndexes

  # GSIのキー属性を収集 (重複排除)
  gsi_attributes = distinct(flatten([
    for gsi in local.gsis : [
      {
        name = gsi.KeyAttributes.PartitionKey.AttributeName
        type = gsi.KeyAttributes.PartitionKey.AttributeType
      },
      {
        name = gsi.KeyAttributes.SortKey.AttributeName
        type = gsi.KeyAttributes.SortKey.AttributeType
      }
    ]
  ]))

  # アイテムのユニークキーを生成（PKとSKの組み合わせ）
  items_map = { for item in local.table_data : "${item.PK.S}_${item.SK.S}" => item }
}

# definition.jsonからDynamoDBテーブルを定義
resource "aws_dynamodb_table" "cms" {
  name           = local.table_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = local.partition_key.AttributeName
  range_key      = local.sort_key.AttributeName

  # パーティションキー
  attribute {
    name = local.partition_key.AttributeName
    type = local.partition_key.AttributeType
  }

  # ソートキー
  attribute {
    name = local.sort_key.AttributeName
    type = local.sort_key.AttributeType
  }

  # GSIのキー属性（動的に生成）
  dynamic "attribute" {
    for_each = { for attr in local.gsi_attributes : attr.name => attr }

    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  # GSI定義（複数のGSIに対応）
  dynamic "global_secondary_index" {
    for_each = { for idx, gsi in local.gsis : gsi.IndexName => gsi }

    content {
      name            = global_secondary_index.key
      hash_key        = global_secondary_index.value.KeyAttributes.PartitionKey.AttributeName
      range_key       = global_secondary_index.value.KeyAttributes.SortKey.AttributeName
      write_capacity  = 5
      read_capacity   = 5
      projection_type = global_secondary_index.value.Projection.ProjectionType
    }
  }

  tags = {
    Name        = "${local.table_name}Table"
    Environment = "dev"
  }
}

# for_eachを使ってすべてのデータを追加
resource "aws_dynamodb_table_item" "cms_items" {
  for_each = local.items_map

  table_name = aws_dynamodb_table.cms.name
  hash_key   = aws_dynamodb_table.cms.hash_key
  range_key  = aws_dynamodb_table.cms.range_key

  item = jsonencode(each.value)
}
