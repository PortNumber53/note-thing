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

		stage('Backend: Test & Migrate') {
			steps {
				withCredentials([
					string(credentialsId: 'prod-database-url-note-thing', variable: 'DATABASE_URL'),
				]) {
					sh '''
						cd backend
						go test ./...
						go run . migrate up
					'''
				}
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
					string(credentialsId: 'prod-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
					string(credentialsId: 'prod-enable-google-oauth-note-thing', variable: 'ENABLE_GOOGLE_OAUTH'),
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
						export VITE_ENABLE_GOOGLE_OAUTH="$ENABLE_GOOGLE_OAUTH"
						npm run build
						npx wrangler secret put BACKEND_URL --config wrangler.jsonc <<EOF
$BACKEND_URL
EOF
						npx wrangler deploy --config wrangler.jsonc
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
					string(credentialsId: 'prod-database-url-note-thing', variable: 'DATABASE_URL'),
					string(credentialsId: 'prod-google-client-id-note-thing', variable: 'GOOGLE_CLIENT_ID'),
					string(credentialsId: 'prod-google-client-secret-note-thing', variable: 'GOOGLE_CLIENT_SECRET'),
					string(credentialsId: 'prod-google-redirect-url-note-thing', variable: 'GOOGLE_REDIRECT_URL'),
					string(credentialsId: 'prod-jwt-secret-note-thing', variable: 'JWT_SECRET'),
					string(credentialsId: 'prod-frontend-url-note-thing', variable: 'FRONTEND_URL'),
					string(credentialsId: 'prod-stripe-secret-key-note-thing', variable: 'STRIPE_SECRET_KEY'),
					string(credentialsId: 'prod-stripe-webhook-secret-note-thing', variable: 'STRIPE_WEBHOOK_SECRET'),
					string(credentialsId: 'prod-admin-emails-note-thing', variable: 'ADMIN_EMAILS'),
				]) {
					sh '''
						export DATABASE_URL="$DATABASE_URL"
						export GOOGLE_CLIENT_ID="$GOOGLE_CLIENT_ID"
						export GOOGLE_CLIENT_SECRET="$GOOGLE_CLIENT_SECRET"
						export GOOGLE_REDIRECT_URL="$GOOGLE_REDIRECT_URL"
						export JWT_SECRET="$JWT_SECRET"
						export FRONTEND_URL="$FRONTEND_URL"
						export STRIPE_SECRET_KEY="$STRIPE_SECRET_KEY"
						export STRIPE_WEBHOOK_SECRET="$STRIPE_WEBHOOK_SECRET"
						export ADMIN_EMAILS="$ADMIN_EMAILS"
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
