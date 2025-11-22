pipeline {
	agent any

	options {
		timestamps()
		disableConcurrentBuilds()
	}

	stages {
		stage('Frontend: Install') {
			steps {
				sh 'cd frontend && npm ci'
			}
		}

		stage('Frontend: Lint') {
			steps {
				sh 'cd frontend && npm run lint'
			}
		}

		stage('Frontend: Build') {
			steps {
				sh 'cd frontend && npm run build'
			}
		}

		stage('Backend: Test') {
			steps {
				sh 'cd backend && go test ./...'
			}
		}

		stage('Deploy Frontend (Cloudflare)') {
			when {
				anyOf {
					branch 'master'
					branch 'main'
				}
			}
			steps {
				withCredentials([
					string(credentialsId: 'cloudflare-api-token', variable: 'CLOUDFLARE_API_TOKEN'),
					string(credentialsId: 'prod-xata-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
				]) {
					sh '''
						cd frontend
						export CF_API_TOKEN="$CLOUDFLARE_API_TOKEN"
						export CLOUDFLARE_API_TOKEN="$CLOUDFLARE_API_TOKEN"
						export DATABASE_URL="$DATABASE_URL"
						export GOOGLE_CLIENT_ID="$GOOGLE_CLIENT_ID"
						export GOOGLE_CLIENT_SECRET="$GOOGLE_CLIENT_SECRET"
						export JWT_SECRET="$JWT_SECRET"
						npm run deploy
					'''
				}
			}
		}

		stage('Deploy Backend') {
			when {
				allOf {
					anyOf {
						branch 'master'
						branch 'main'
					}
					expression { return env.BACKEND_DEPLOY_ENABLED == 'true' }
				}
			}
			steps {
				withCredentials([
					string(credentialsId: 'prod-xata-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
				]) {
					sh '''
						echo "Backend deploy enabled. Running BACKEND_DEPLOY_COMMAND..."
						if [ -z "$BACKEND_DEPLOY_COMMAND" ]; then
							echo "BACKEND_DEPLOY_COMMAND is not set"
							exit 1
						fi
						export DATABASE_URL="$DATABASE_URL"
						export GOOGLE_CLIENT_ID="$GOOGLE_CLIENT_ID"
						export GOOGLE_CLIENT_SECRET="$GOOGLE_CLIENT_SECRET"
						export JWT_SECRET="$JWT_SECRET"
						eval "$BACKEND_DEPLOY_COMMAND"
					'''
				}
			}
		}
	}

	post {
		always {
			archiveArtifacts artifacts: 'frontend/dist/**', allowEmptyArchive: true
		}
	}
}
