pipeline {
	agent any

	options {
		ansiColor('xterm')
		timestamps()
		disableConcurrentBuilds()
	}

	environment {
		NODE_ENV = 'production'
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
				withCredentials([string(credentialsId: 'cloudflare-api-token', variable: 'CLOUDFLARE_API_TOKEN')]) {
					sh '''
						cd frontend
						export CF_API_TOKEN="$CLOUDFLARE_API_TOKEN"
						export CLOUDFLARE_API_TOKEN="$CLOUDFLARE_API_TOKEN"
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
				sh '''
					echo "Backend deploy enabled. Running BACKEND_DEPLOY_COMMAND..."
					if [ -z "$BACKEND_DEPLOY_COMMAND" ]; then
						echo "BACKEND_DEPLOY_COMMAND is not set"
						exit 1
					fi
					eval "$BACKEND_DEPLOY_COMMAND"
				'''
			}
		}
	}

	post {
		always {
			archiveArtifacts artifacts: 'frontend/dist/**', allowEmptyArchive: true
		}
	}
}
