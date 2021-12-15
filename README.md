# application-backend
## Tech Stack
1. `Gin`:golang web 框架，支援middleware、crash-free、JSON validation、Routes grouping、gin.Context。
2. `Docker`：協助建置、測試並且打包成一個獨立網頁應用，只需建置一次 Docker image，即可在任意機器、平台、服務上執行，避免環境設定、套件安裝等繁雜工作。
3. `Google Cloud Platform`：協助我們整合 CI/CD 自動進行建置、測試、部署至雲端，讓我們有更多餘力專注在開發上。
4. `Github Actions`:用於在GitHub測試、封裝、發佈或部署任何專案，並可藉此建置 CI/CD 功能。

## Development
### Docker
主要提供給前端開發人員使用：

- `docker-compose pull`：pull backend server image，名稱為 `blawhi2435/whispering-corner`
- `docker-compose up -d`：啟動 API Server、MySQL、Redis，Listen [http://localhost:6969](http://localhost:6969)。
- `docker-compose down`：停止 API Server、MySQL、Redis。

### GCP部署
使用service account讓Github Actions將專案部署至Cloud Run上。
- docker image path: `asia-east1-docker.pkg.dev/nth-weft-328504/backend-repo/application-backend`
- cloud run service name: `application-backend`
- cloud run service url: [ https://application-backend-mhkzpmkvca-de.a.run.app](https://application-backend-mhkzpmkvca-de.a.run.app)
- cloud run url access permission: `allUsers`

### Github Actions
#### Workflow
目前伺服器是測試即正式環境，所以暫時使用2種workflow分別處理CI及CD，目前也尚未使用auto-merge的action。
- build_and_test: 在push至main及pull request時觸發建置及測試。
- release: 當release一個版本時將會自動部署至Google Cloud Run。

#### Rules
Container image name 要遵守Artifact Registry命名規則，否則在Cloud Build的步驟失敗。
- Container image names format: `LOCATION-docker.pkg.dev/PROJECT-ID/REPOSITORY/IMAGE:TAG`

#### ENV
Workflow中使用到的環境參數介紹。
- PROJECT_ID: Google Project ID。
- REPO_NAME: Artifact Registry存放區上的名稱。
- IMAGE_NAME: Artifact Registry上的映像名稱。
- TAG: image version，目前都先使用release。

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