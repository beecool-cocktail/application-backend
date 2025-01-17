# application-backend
## Tech Stack
1. `Gin`:golang web 框架，支援middleware、crash-free、JSON validation、Routes grouping、gin.Context。
2. `Gorm`:golang ORM 框架，支援MySQL、PostgreSQL、SQLite、SQL Server等資料庫，提供開發者更方便操作資料庫。
3. `Docker`：協助建置、測試並且打包成一個獨立網頁應用，只需建置一次 Docker image，即可在任意機器、平台、服務上執行，避免環境設定、套件安裝等繁雜工作。
4. `Google Cloud Platform`：協助整合 CI/CD 自動進行建置、測試、部署至雲端，有更多餘力專注在開發上。
5. `Github Actions`:用於在GitHub測試、封裝、發佈或部署任何專案，並可藉此建置 CI/CD 功能。

## Development
### Docker
主要提供給前端開發人員使用：

直接使用寫好的 script 即刻啟動整個後端服務，不需安裝 MySQL、Redis、Elastic、Golang 等環境。
- `./script/start.sh runall`: 啟動完整的後端服務。
- `./script/start.sh runsimplify`: 啟動簡化的後端服務(不啟動 Elastic)。
- `./script/start.sh down`: 停止所有後端服務。

### Github Actions
#### Workflow
目前伺服器是測試即正式環境，所以暫時使用2種workflow分別處理CI及CD，目前也尚未使用auto-merge的action。
- build_and_test: 在push及pull request時觸發建置及測試。
- update_docker_hub: 在 push 至 master 或 beta 版本會同時更新至Docker Hub。
- release: 當release一個版本時將會自動部署至Google Compute Engine。

### Go-Swagger
#### Installing 
Homebrew/Linuxbrew
```
brew tap go-swagger/go-swagger
brew install go-swagger
```
 
 > Note: <br />
 > 其他安裝方式可以參考[這裡](https://goswagger.io/install.html)。

### Run Swagger
```bash
# start swagger server
$ swagger serve -F=swagger swagger.yaml

# create api document markdown
$ swagger generate markdown -f ./swagger.yaml --output swagger.md
```

## Reference
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Building Go Web Applications and Micro services Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin)
- [Artifact Registry](https://cloud.google.com/artifact-registry)
- [How to Deploy Static Site to GCP Cloud Run](https://galtz.netlify.app/gcp-static-site/)
- [Deploy To Google Cloud Run Using Github Actions](https://towardsdatascience.com/deploy-to-google-cloud-run-using-github-actions-590ecf957af0)