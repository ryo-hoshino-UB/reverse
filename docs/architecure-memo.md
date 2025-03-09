## router

- リクエストを受け取り、ドメインに変換・バリデーションし service に渡す
- gateway のインスタンスを作成して service に DI

## service

- router から渡された値に加えて必要なデータを DB から取得し、ビジネスロジックを呼び出して実行する
- service のコンストラクタの引数には repository の interface が設定されている
- つまり service のメソッドは、repository のメソッドを呼び出す

## domain

- service 層から呼び出されるビジネスロジックの実装
