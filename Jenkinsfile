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
					string(credentialsId: 'prod-backend-url-note-thing', variable: 'BACKEND_URL'),
					string(credentialsId: 'prod-xata-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
				]) {
					sh '''
						cd frontend
						export CF_API_TOKEN="$CLOUDFLARE_API_TOKEN"
						export CLOUDFLARE_API_TOKEN="$CLOUDFLARE_API_TOKEN"
						export BACKEND_URL="$BACKEND_URL"
						export DATABASE_URL="$DATABASE_URL"
						export GOOGLE_CLIENT_ID="$GOOGLE_CLIENT_ID"
						export GOOGLE_CLIENT_SECRET="$GOOGLE_CLIENT_SECRET"
						export JWT_SECRET="$JWT_SECRET"
						npm run build
						npx wrangler deploy --var BACKEND_URL="$BACKEND_URL"
					'''
				}
			}
		}

		stage('Deploy Backend') {
			when {
				anyOf {
					branch 'master'
					branch 'main'
				}
			}
			steps {
				withCredentials([
					sshUserPrivateKey(credentialsId: 'Jenkins-private-key', keyFileVariable: 'SSH_KEY'),
					string(credentialsId: 'prod-xata-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
				]) {
					sh '''
						export DATABASE_URL="$DATABASE_URL"
						export GOOGLE_CLIENT_ID="$GOOGLE_CLIENT_ID"
						export GOOGLE_CLIENT_SECRET="$GOOGLE_CLIENT_SECRET"
						export JWT_SECRET="$JWT_SECRET"
						cd backend
						chmod +x deploy/deploy_backend.sh
						./deploy/deploy_backend.sh "$SSH_KEY"
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
