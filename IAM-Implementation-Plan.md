# IAM System Implementation Plan
# Generic Multi-Tenant Identity and Access Management System

---

## Document Control

| Attribute | Value |
|-----------|-------|
| **Document Title** | IAM System Implementation Plan |
| **Version** | 2.0 |
| **Created Date** | 2026-01-22 |
| **Last Updated** | 2026-02-03 |
| **Project Duration** | 18 weeks (4.5 months) |
| **Team Size** | 5-7 developers |

### Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-01-22 | Initial implementation plan |
| 2.0 | 2026-02-03 | Major revision: Added Masterdata integration tasks, Email service setup, Token blacklist implementation, API documentation moved earlier, Compliance verification tasks, Bulk import clarification (deferred), Corrected migration file count (16), Added Security Engineer role, Extended timeline to 18 weeks |

---

## Executive Summary

This implementation plan outlines the development strategy for the Generic Multi-Tenant IAM System. The project is organized into 4 phases spanning 18 weeks, with clear milestones, deliverables, and dependencies.

**Key Dates:**
- Project Start: Week 1 (Q1 2026)
- API Contract Draft: Week 4
- Alpha Release: Week 9
- Beta Release: Week 14
- Production Release: Week 18

**Key Changes from v1.0:**
- Extended timeline from 16 to 18 weeks for quality assurance
- Added dedicated tasks for Masterdata service integration
- Added Email service (SMTP) setup and configuration
- Added explicit Token blacklist implementation
- Moved API documentation (OpenAPI/Postman) to Phase 1
- Added comprehensive compliance verification tasks
- Clarified database migration count (16 files per TRD)
- Added Security Engineer role to team structure
- Deferred Bulk User Import to Phase 2 (future release)

---

## 1. Project Overview

### 1.1 Objectives

1. Deliver a production-ready, multi-tenant IAM system
2. Achieve 99.9% uptime SLA capability
3. Pass security audit (OWASP Top 10, ISO 27001 alignment)
4. Enable < 2 day integration time for new applications
5. Achieve full compliance with ISO 27001 and PCI DSS (where applicable)

### 1.2 Success Criteria

| Criteria | Target | Measurement |
|----------|--------|-------------|
| All PRD features implemented | 100% | Feature checklist |
| Unit test coverage | ≥ 80% | Code coverage tool |
| Integration test coverage | ≥ 70% | Test report |
| Security vulnerabilities | 0 critical, 0 high | Security scan |
| API response time (p95) | < 500ms | Load testing |
| Documentation complete | 100% | Doc review |
| Compliance checklist | 100% pass | Audit checklist |
| API contract validation | 100% | Postman tests |

### 1.3 Team Structure

| Role | Count | Responsibilities |
|------|-------|------------------|
| Tech Lead | 1 | Architecture, code review, technical decisions |
| Backend Developer | 2-3 | API development, business logic |
| DevOps Engineer | 1 | Infrastructure, CI/CD, deployment |
| QA Engineer | 1 | Testing, quality assurance |
| Security Engineer | 0.5 (shared) | Security review, penetration testing, compliance |
| Product Owner | 1 (part-time) | Requirements, acceptance |

> **Note:** Security Engineer is a shared resource from the security team, allocated 50% to this project during Phases 3-4.

### 1.4 Out of Scope (Deferred to Future Releases)

| Item | Rationale | Target Release |
|------|-----------|----------------|
| Bulk User Import (CSV/Excel) | Complexity; prioritize core features | v1.1 |
| Resource-level ACL | RBAC with branch scoping covers v1.0 | v2.0 |
| SSO/SAML Integration | Enterprise feature | v2.0 |
| TOTP Authenticator Apps | Alternative MFA | v1.1 |
| Webhooks/Event Notifications | Apps query IAM for now | v1.1 |

---

## 2. Phase Overview

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    PROJECT TIMELINE (18 WEEKS)                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────┘

Week:  1   2   3   4   5   6   7   8   9  10  11  12  13  14  15  16  17  18
       │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
       ├───┴───┴───┴───┤   │   │   │   │   │   │   │   │   │   │   │   │   │
       │   PHASE 1     │   │   │   │   │   │   │   │   │   │   │   │   │   │
       │  Foundation   │   │   │   │   │   │   │   │   │   │   │   │   │   │
       │  + API Docs   │   │   │   │   │   │   │   │   │   │   │   │   │   │
       │  (4 weeks)    │   │   │   │   │   │   │   │   │   │   │   │   │   │
       └───────────────┴───┴───┴───┴───┤   │   │   │   │   │   │   │   │   │
                       │   PHASE 2     │   │   │   │   │   │   │   │   │   │
                       │ Core Features │   │   │   │   │   │   │   │   │   │
                       │  (5 weeks)    │   │   │   │   │   │   │   │   │   │
                       └───────────────┴───┴───┴───┴───┤   │   │   │   │   │
                                       │   PHASE 3     │   │   │   │   │   │
                                       │   Advanced    │   │   │   │   │   │
                                       │  (5 weeks)    │   │   │   │   │   │
                                       └───────────────┴───┴───┴───┴───┤   │
                                                       │   PHASE 4     │   │
                                                       │  Production   │   │
                                                       │  (4 weeks)    │   │
                                                       └───────────────┴───┘

Milestones:
  ▲ Week 4: Infrastructure Ready + API Contract Draft
  ▲ Week 9: Alpha Release (Internal Testing)
  ▲ Week 14: Beta Release (Staging)
  ▲ Week 18: Production Release
```

---

## 3. Phase 1: Foundation (Weeks 1-4)

### 3.1 Objectives

- Set up development infrastructure
- Establish coding standards and patterns
- Implement database schema and migrations
- Create core service architecture
- **NEW:** Draft API contract (OpenAPI + Postman collection)
- **NEW:** Set up Email service infrastructure
- **NEW:** Integrate Masterdata service (stub if unavailable)

### 3.2 Sprint 1 (Weeks 1-2): Project Setup

#### 3.2.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 1.1.1 | Initialize Go project with module structure | Backend | 4 | P0 |
| 1.1.2 | Set up Docker Compose for local dev | DevOps | 8 | P0 |
| 1.1.3 | Configure PostgreSQL 18 with extensions (uuidv7) | DevOps | 4 | P0 |
| 1.1.4 | Configure Redis 7 | DevOps | 2 | P0 |
| 1.1.5 | Configure OpenSearch cluster | DevOps | 8 | P0 |
| 1.1.6 | Set up HashiCorp Vault (dev mode) | DevOps | 4 | P0 |
| 1.1.7 | Create database migration framework (golang-migrate) | Backend | 8 | P0 |
| 1.1.8 | Implement all 16 database migrations (per TRD Section 3) | Backend | 16 | P0 |
| 1.1.9 | Set up CI pipeline (lint, test, build) | DevOps | 8 | P0 |
| 1.1.10 | Configure golangci-lint rules | Tech Lead | 4 | P1 |
| 1.1.11 | Create Makefile with common commands | Backend | 4 | P1 |
| 1.1.12 | Write seed data scripts (platform admin, test tenant) | Backend | 8 | P1 |
| **1.1.13** | **Set up SMTP service (Mailhog for dev, config for prod)** | DevOps | 4 | P0 |
| **1.1.14** | **Create email templates (OTP, welcome, password reset)** | Backend | 8 | P1 |
| **1.1.15** | **Set up Masterdata service stub (mock API)** | Backend | 8 | P1 |

#### 3.2.2 Deliverables

- [ ] Working local development environment
- [ ] All 16 database migration files (per TRD-Section3-DatabaseDesign.md)
- [ ] CI pipeline (GitHub Actions)
- [ ] Seed data for development/testing
- [ ] **Email service configured (Mailhog for local, SMTP config for staging/prod)**
- [ ] **Masterdata service stub with mock data (gender, marital status, etc.)**

#### 3.2.3 Definition of Done

- All services start with `make dev`
- Migrations run successfully (16 files)
- Linting passes with no errors
- README with setup instructions
- Email sending works locally (Mailhog captures)
- Masterdata mock returns valid responses

### 3.3 Sprint 2 (Weeks 3-4): Core Architecture + API Contract

#### 3.3.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 1.2.1 | Implement domain models (entities) | Backend | 16 | P0 |
| 1.2.2 | Create repository interfaces | Backend | 8 | P0 |
| 1.2.3 | Implement PostgreSQL repositories | Backend | 24 | P0 |
| 1.2.4 | Implement Redis repository (sessions, rate limits, OTP) | Backend | 8 | P0 |
| 1.2.5 | Create Vault client wrapper | Backend | 8 | P0 |
| 1.2.6 | Implement configuration loading (env + Vault) | Backend | 8 | P0 |
| 1.2.7 | Set up Fiber HTTP framework | Backend | 4 | P0 |
| 1.2.8 | Create middleware stack (logging, recovery, CORS) | Backend | 16 | P0 |
| 1.2.9 | Implement structured logging (Zap) | Backend | 4 | P1 |
| 1.2.10 | Create error handling framework | Backend | 8 | P0 |
| 1.2.11 | Write unit tests for repositories | QA | 16 | P1 |
| 1.2.12 | Create integration test framework | QA | 8 | P1 |
| **1.2.13** | **Create OpenAPI specification (draft, all endpoints)** | Tech Lead | 16 | P0 |
| **1.2.14** | **Create Postman collection (all endpoints with examples)** | Tech Lead | 12 | P0 |
| **1.2.15** | **Implement Masterdata client (API integration + Redis cache)** | Backend | 12 | P0 |
| **1.2.16** | **Implement Email service wrapper (template rendering + SMTP)** | Backend | 8 | P0 |

#### 3.3.2 Deliverables

- [ ] Complete domain model layer
- [ ] Repository implementations (PostgreSQL, Redis)
- [ ] Vault integration for secrets
- [ ] HTTP server with middleware chain
- [ ] Logging and error handling
- [ ] **OpenAPI specification (draft) - openapi-iam-v1.yaml**
- [ ] **Postman collection with all endpoints**
- [ ] **Masterdata client with Redis caching (1 hour TTL)**
- [ ] **Email service with templating**

#### 3.3.3 Definition of Done

- Repository tests pass (≥80% coverage)
- Application starts and connects to all services
- Health check endpoint works
- Code review completed
- **OpenAPI spec validated (no errors)**
- **Postman collection runs against mock server**
- **Masterdata client handles service unavailability gracefully**

### 3.4 Phase 1 Milestone: Infrastructure Ready + API Contract Draft

**Checkpoint Criteria:**
- [ ] All infrastructure components running
- [ ] Database schema deployed (16 migrations)
- [ ] CI/CD pipeline functional
- [ ] Core architecture in place
- [ ] Development workflow documented
- [ ] **OpenAPI specification (draft) available for frontend teams**
- [ ] **Postman collection available for integration testing**
- [ ] **Email service functional**
- [ ] **Masterdata integration working (with fallback)**

---

## 4. Phase 2: Core Features (Weeks 5-9)

### 4.1 Objectives

- Implement authentication (login, registration, PIN)
- Implement user management
- Create authorization foundation
- Achieve internal testing readiness
- **NEW:** Implement token blacklist for secure logout

### 4.2 Sprint 3 (Weeks 5-6): Authentication

#### 4.2.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 2.1.1 | Implement password hashing (bcrypt, cost 12) | Backend | 4 | P0 |
| 2.1.2 | Implement PIN hashing/validation (bcrypt, cost 12) | Backend | 4 | P0 |
| 2.1.3 | Create JWT service (RS256, key rotation) | Backend | 16 | P0 |
| 2.1.4 | Implement JWKS endpoint (/.well-known/jwks.json) | Backend | 4 | P0 |
| 2.1.5 | Implement login endpoint (POST /auth/login) | Backend | 16 | P0 |
| 2.1.6 | Implement PIN verification endpoint (POST /auth/verify-pin) | Backend | 8 | P0 |
| 2.1.7 | Implement token refresh endpoint (POST /auth/refresh) | Backend | 8 | P0 |
| 2.1.8 | Implement logout endpoint (POST /auth/logout) | Backend | 4 | P0 |
| **2.1.9** | **Implement token blacklist in Redis (for logout)** | Backend | 8 | P0 |
| 2.1.10 | Implement account lockout logic (5 failed attempts) | Backend | 8 | P0 |
| 2.1.11 | Create auth middleware (JWT validation + blacklist check) | Backend | 12 | P0 |
| 2.1.12 | Implement rate limiting (per IP, per tenant) | Backend | 16 | P0 |
| 2.1.13 | Write authentication tests | QA | 16 | P0 |
| 2.1.14 | Implement Google OAuth (POST /auth/oauth/google) | Backend | 12 | P1 |
| **2.1.15** | **Implement Google OAuth callback (POST /auth/oauth/google/callback)** | Backend | 8 | P1 |

#### 4.2.2 Deliverables

- [ ] POST /auth/login
- [ ] POST /auth/verify-pin
- [ ] POST /auth/refresh
- [ ] POST /auth/logout
- [ ] GET /.well-known/jwks.json
- [ ] Rate limiting middleware
- [ ] Auth middleware (with token blacklist check)
- [ ] **Token blacklist in Redis (TTL = token expiry)**
- [ ] POST /auth/oauth/google
- [ ] POST /auth/oauth/google/callback

#### 4.2.3 Definition of Done

- All auth endpoints functional
- Rate limiting working (429 responses with headers)
- Token rotation working
- Account lockout tested (locks after 5 failures)
- **Logout invalidates token immediately (blacklist)**
- 80%+ test coverage
- Postman collection updated with auth tests

### 4.3 Sprint 4 (Weeks 7-8): User Management & Registration

#### 4.3.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 2.2.1 | Implement user registration (POST /auth/register) | Backend | 16 | P0 |
| 2.2.2 | Implement email OTP generation/validation (Redis) | Backend | 8 | P0 |
| 2.2.3 | Implement email verification (POST /auth/verify-email) | Backend | 8 | P0 |
| 2.2.4 | Implement PIN setup (POST /auth/setup-pin) | Backend | 8 | P0 |
| 2.2.5 | Implement password change (POST /auth/change-password) | Backend | 8 | P0 |
| 2.2.6 | Implement forgot/reset password flow | Backend | 16 | P0 |
| **2.2.7** | **Implement password history check (prevent reuse of last 5)** | Backend | 4 | P0 |
| **2.2.8** | **Implement password history pruning (auto-cleanup)** | Backend | 4 | P1 |
| 2.2.9 | Implement user CRUD endpoints | Backend | 16 | P0 |
| 2.2.10 | Implement user approval workflow | Backend | 8 | P0 |
| 2.2.11 | Implement user unlock (POST /users/{id}/unlock) | Backend | 4 | P0 |
| 2.2.12 | Implement user profile endpoint (GET/PUT /users/me/profile) | Backend | 8 | P0 |
| **2.2.13** | **Validate profile fields via Masterdata (gender, marital status)** | Backend | 4 | P0 |
| 2.2.14 | Write user management tests | QA | 16 | P0 |
| 2.2.15 | Integration testing (full registration flow) | QA | 16 | P0 |

#### 4.3.2 Deliverables

- [ ] POST /auth/register
- [ ] POST /auth/verify-email
- [ ] POST /auth/setup-pin
- [ ] POST /auth/change-password
- [ ] POST /auth/forgot-password
- [ ] POST /auth/reset-password
- [ ] User CRUD endpoints (POST, GET, PUT, DELETE /users)
- [ ] User approval workflow (POST /users/{id}/approve, /reject)
- [ ] **Password history enforcement (last 5 passwords)**
- [ ] **Masterdata validation for profile fields**

#### 4.3.3 Definition of Done

- Complete user lifecycle working
- Email sending functional (OTP, welcome, reset)
- Password policy enforced (complexity + history)
- Integration tests pass
- Postman collection updated

### 4.4 Sprint 5 (Week 9): Alpha Stabilization

#### 4.4.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 2.3.1 | Bug fixes from Sprint 3-4 | Backend | 24 | P0 |
| 2.3.2 | Performance optimization (login flow) | Backend | 8 | P1 |
| 2.3.3 | API documentation update (OpenAPI + Postman) | Tech Lead | 8 | P0 |
| 2.3.4 | Internal testing (team dogfooding) | All | 16 | P0 |
| 2.3.5 | Security self-assessment (OWASP checklist) | Backend | 8 | P0 |
| 2.3.6 | Alpha release preparation | DevOps | 8 | P0 |

#### 4.4.2 Deliverables

- [ ] All critical bugs fixed
- [ ] Alpha environment deployed
- [ ] Updated API documentation
- [ ] Security self-assessment report

### 4.5 Phase 2 Milestone: Alpha Release

**Checkpoint Criteria:**
- [ ] Authentication fully functional
- [ ] User registration/management working
- [ ] Basic authorization in place
- [ ] Internal team testing begins
- [ ] Test coverage ≥ 70%
- [ ] **Token blacklist working**
- [ ] **Password history enforced**
- [ ] **Masterdata integration validated**
- [ ] **OpenAPI spec updated to match implementation**

---

## 5. Phase 3: Advanced Features (Weeks 10-14)

### 5.1 Objectives

- Complete RBAC implementation
- Implement organization management
- Set up audit logging
- Achieve beta readiness
- **NEW:** Begin compliance verification

### 5.2 Sprint 6 (Weeks 10-11): Authorization (RBAC)

#### 5.2.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 3.1.1 | Implement application CRUD | Backend | 8 | P0 |
| 3.1.2 | Implement role CRUD | Backend | 16 | P0 |
| 3.1.3 | Implement permission CRUD | Backend | 12 | P0 |
| 3.1.4 | Implement role-permission assignment | Backend | 8 | P0 |
| 3.1.5 | Implement user-role assignment | Backend | 16 | P0 |
| 3.1.6 | Implement permission check service | Backend | 16 | P0 |
| 3.1.7 | Create permission middleware | Backend | 8 | P0 |
| 3.1.8 | Implement branch-scoped roles | Backend | 12 | P1 |
| 3.1.9 | Implement permission caching (Redis) | Backend | 8 | P0 |
| 3.1.10 | Write authorization tests | QA | 16 | P0 |
| 3.1.11 | Performance testing (permission check < 100ms) | QA | 8 | P1 |

#### 5.2.2 Deliverables

- [ ] Application management endpoints
- [ ] Role management endpoints
- [ ] Permission management endpoints
- [ ] User role assignment endpoints
- [ ] Permission check API (POST /permissions/check)
- [ ] Permission caching (Redis, 5 min TTL)

#### 5.2.3 Definition of Done

- All RBAC endpoints functional
- Permission check < 100ms (p95)
- Branch-scoped permissions working
- Comprehensive test coverage
- Postman collection updated

### 5.3 Sprint 7 (Weeks 12-13): Organization & Audit

#### 5.3.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 3.2.1 | Implement tenant management | Backend | 16 | P0 |
| 3.2.2 | Implement branch CRUD | Backend | 12 | P0 |
| 3.2.3 | Implement user-branch assignment | Backend | 8 | P0 |
| 3.2.4 | Create OpenSearch client | Backend | 8 | P0 |
| 3.2.5 | Implement audit service (async logging) | Backend | 16 | P0 |
| 3.2.6 | Create audit event types (per PRD F5) | Backend | 4 | P0 |
| 3.2.7 | Implement audit log query API | Backend | 12 | P0 |
| 3.2.8 | Create OpenSearch index templates | DevOps | 4 | P0 |
| 3.2.9 | Implement ILM policies (90 days hot, 7 years archive) | DevOps | 4 | P0 |
| 3.2.10 | Write organization tests | QA | 12 | P0 |
| 3.2.11 | End-to-end testing | QA | 16 | P0 |

#### 5.3.2 Deliverables

- [ ] Tenant management (platform admin only)
- [ ] Branch management endpoints
- [ ] Audit logging to OpenSearch
- [ ] Audit query API (GET /audit/logs)

#### 5.3.3 Definition of Done

- All organization features working
- Audit events logged for all operations
- Audit log retention configured
- E2E tests pass

### 5.4 Sprint 8 (Week 14): Beta Preparation & Staging

#### 5.4.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 3.3.1 | Staging environment setup | DevOps | 16 | P0 |
| 3.3.2 | Staging Vault configuration | DevOps | 8 | P0 |
| 3.3.3 | Staging database setup (with test data) | DevOps | 8 | P0 |
| 3.3.4 | Staging OpenSearch setup | DevOps | 4 | P0 |
| 3.3.5 | Beta release deployment | DevOps | 8 | P0 |
| **3.3.6** | **Compliance pre-assessment (ISO 27001 checklist)** | Security | 16 | P0 |
| **3.3.7** | **Security scan (Trivy, gosec)** | Security | 8 | P0 |
| 3.3.8 | API documentation finalization | Tech Lead | 8 | P0 |
| 3.3.9 | Postman collection finalization (all tests passing) | Tech Lead | 8 | P0 |

#### 5.4.2 Deliverables

- [ ] Staging environment operational
- [ ] Beta deployment successful
- [ ] **Compliance pre-assessment report**
- [ ] **Security scan report (0 critical/high)**
- [ ] Finalized API documentation

### 5.5 Phase 3 Milestone: Beta Release

**Checkpoint Criteria:**
- [ ] All PRD features implemented
- [ ] Staging environment operational
- [ ] Security scan passed (0 critical/high)
- [ ] Performance targets met
- [ ] Beta testing initiated
- [ ] **Compliance pre-assessment completed**
- [ ] **API documentation finalized (OpenAPI + Postman)**

---

## 6. Phase 4: Production Readiness (Weeks 15-18)

### 6.1 Objectives

- Harden security
- Optimize performance
- Complete documentation
- **NEW:** Full compliance verification
- Production deployment

### 6.2 Sprint 9 (Weeks 15-16): Security & Performance

#### 6.2.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 4.1.1 | Security audit preparation | Security | 8 | P0 |
| 4.1.2 | Implement security headers (HSTS, CSP, etc.) | Backend | 4 | P0 |
| 4.1.3 | CORS configuration (production whitelist) | Backend | 4 | P0 |
| 4.1.4 | Input validation hardening | Backend | 8 | P0 |
| 4.1.5 | Performance profiling (pprof) | Backend | 16 | P0 |
| 4.1.6 | Query optimization (EXPLAIN ANALYZE) | Backend | 16 | P0 |
| 4.1.7 | Cache optimization (Redis hit rate > 90%) | Backend | 8 | P1 |
| 4.1.8 | Load testing (k6, 1000 concurrent users) | QA | 16 | P0 |
| **4.1.9** | **Penetration testing** | Security | 24 | P0 |
| 4.1.10 | Bug fixes from beta | Backend | 24 | P0 |
| **4.1.11** | **Compliance verification (ISO 27001 checklist)** | Security | 16 | P0 |
| **4.1.12** | **PCI DSS assessment (if applicable)** | Security | 8 | P1 |

#### 6.2.2 Deliverables

- [ ] Security audit report (clean)
- [ ] Performance benchmarks (all targets met)
- [ ] Load test results (1000 users, < 500ms p95)
- [ ] Bug fixes implemented
- [ ] **Penetration test report**
- [ ] **Compliance verification report**

#### 6.2.3 Definition of Done

- No critical/high vulnerabilities
- Performance targets achieved
- All beta feedback addressed
- Compliance checklist 100% pass

### 6.3 Sprint 10 (Weeks 17-18): Documentation & Launch

#### 6.3.1 Tasks

| ID | Task | Owner | Est. Hours | Priority |
|----|------|-------|------------|----------|
| 4.2.1 | OpenAPI specification (final) | Backend | 8 | P0 |
| 4.2.2 | Postman collection (final, with environments) | Backend | 8 | P0 |
| 4.2.3 | API documentation (developer portal) | Backend | 16 | P0 |
| 4.2.4 | Integration guide (for consuming applications) | Tech Lead | 16 | P0 |
| 4.2.5 | Operations runbook | DevOps | 16 | P0 |
| 4.2.6 | Production Kubernetes manifests | DevOps | 16 | P0 |
| 4.2.7 | Production deployment | DevOps | 16 | P0 |
| 4.2.8 | Monitoring dashboards (Grafana) | DevOps | 8 | P0 |
| 4.2.9 | Alerting rules (PagerDuty/Opsgenie) | DevOps | 8 | P0 |
| 4.2.10 | Backup/restore procedures (tested) | DevOps | 8 | P0 |
| 4.2.11 | Final regression testing | QA | 16 | P0 |
| 4.2.12 | User acceptance testing | Product | 16 | P0 |
| 4.2.13 | Go-live support | All | 16 | P0 |
| **4.2.14** | **Compliance documentation package** | Security | 8 | P0 |

#### 6.3.2 Deliverables

- [ ] Complete API documentation (OpenAPI + Postman)
- [ ] Integration guide for developers
- [ ] Operations runbook
- [ ] Production deployment
- [ ] Monitoring & alerting
- [ ] **Compliance documentation package**

#### 6.3.3 Definition of Done

- All documentation complete
- Production deployment successful
- Monitoring operational
- UAT passed
- Compliance package approved

### 6.4 Phase 4 Milestone: Production Release

**Checkpoint Criteria:**
- [ ] Production deployment complete
- [ ] Documentation published
- [ ] Monitoring active
- [ ] Support procedures in place
- [ ] First customer onboarded
- [ ] **Compliance certification (if required)**
- [ ] **Security sign-off obtained**

---

## 7. Risk Management

### 7.1 Identified Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Resource availability | High | Medium | Cross-train team members |
| Scope creep | Medium | High | Strict change control |
| Security vulnerabilities | High | Medium | Regular security scans, pen testing |
| Performance issues | High | Low | Early load testing |
| Integration complexity | Medium | Medium | Mock external services |
| Infrastructure issues | High | Low | IaC, disaster recovery |
| **Masterdata service unavailable** | Medium | Medium | **Redis caching + graceful degradation** |
| **Email delivery failures** | Medium | Low | **Retry queue + monitoring** |
| **Compliance gaps** | High | Medium | **Early assessment + dedicated Security Engineer** |

### 7.2 Contingency Plans

| Scenario | Response |
|----------|----------|
| Key person unavailable | Documentation + pair programming |
| Critical bug in production | Rollback procedure + hotfix process |
| Security breach | Incident response plan + communication |
| Missed deadline | Scope reduction to MVP |
| **Masterdata service down** | Use cached values (stale OK for reference data) |
| **SMTP failure** | Queue emails + alert ops team |
| **Compliance audit failure** | Remediation sprint + re-audit |

---

## 8. Dependencies

### 8.1 External Dependencies

| Dependency | Owner | Status | Risk | Mitigation |
|------------|-------|--------|------|------------|
| Google OAuth credentials | Platform team | Pending | Low | Fallback to email/password |
| SMTP service | Infrastructure | Available | Low | Mailhog for dev/test |
| Domain + SSL certificates | Infrastructure | Pending | Low | Can use self-signed for staging |
| Production Kubernetes cluster | Infrastructure | In progress | Medium | Local Docker fallback |
| Vault production instance | Infrastructure | In progress | Medium | Env vars fallback |
| **Masterdata Service** | Platform team | **Stub available** | Medium | **Use stub until ready** |

### 8.2 Internal Dependencies

| Phase | Depends On |
|-------|------------|
| Phase 2 | Phase 1 completion |
| Phase 3 | Phase 2 (auth + users) |
| Phase 4 | Phase 3 completion |
| Production | Staging validation |
| **RBAC** | **User management complete** |
| **Audit logging** | **All feature modules emit events** |

---

## 9. Resource Allocation

### 9.1 Sprint Capacity (per sprint = 2 weeks, except Sprint 5 & 8)

| Role | Hours/Week | Per Sprint (2 weeks) |
|------|------------|----------------------|
| Backend Developer (x2) | 40 each | 160 |
| Tech Lead | 30 | 60 |
| DevOps Engineer | 40 | 80 |
| QA Engineer | 40 | 80 |
| Security Engineer (0.5) | 20 | 40 |
| **Total** | | **420** |

### 9.2 Effort Distribution by Phase

| Phase | Backend | DevOps | QA | Security | Total |
|-------|---------|--------|----|----|-----|
| Phase 1 | 50% | 30% | 10% | 10% | 100% |
| Phase 2 | 60% | 10% | 25% | 5% | 100% |
| Phase 3 | 50% | 20% | 20% | 10% | 100% |
| Phase 4 | 35% | 25% | 20% | 20% | 100% |

---

## 10. Quality Gates

### 10.1 Code Quality

| Metric | Threshold | Tool |
|--------|-----------|------|
| Test coverage | ≥ 80% | go test -cover |
| Lint issues | 0 errors | golangci-lint |
| Code duplication | < 5% | dupl |
| Cyclomatic complexity | < 15 | gocyclo |

### 10.2 Security Gates

| Check | Threshold | Tool |
|-------|-----------|------|
| Critical vulnerabilities | 0 | Trivy, gosec |
| High vulnerabilities | 0 | Trivy, gosec |
| OWASP compliance | Pass | Manual review |
| Dependency vulnerabilities | 0 critical | go mod audit |
| **Penetration test** | Pass | External auditor |

### 10.3 Performance Gates

| Metric | Target | Tool |
|--------|--------|------|
| Login p95 | < 500ms | k6 |
| Permission check p95 | < 100ms | k6 |
| Concurrent users | 1000+ | k6 |
| Error rate under load | < 0.1% | k6 |

### 10.4 API Contract Gates

| Check | Threshold | Tool |
|-------|-----------|------|
| OpenAPI validation | 0 errors | Spectral |
| Postman tests | 100% pass | Newman |
| Response schema validation | 100% match | Postman |

---

## 11. Communication Plan

### 11.1 Meetings

| Meeting | Frequency | Participants | Purpose |
|---------|-----------|--------------|---------|
| Daily Standup | Daily | All team | Progress, blockers |
| Sprint Planning | Bi-weekly | All team | Sprint goals |
| Sprint Review | Bi-weekly | All + stakeholders | Demo, feedback |
| Sprint Retro | Bi-weekly | All team | Process improvement |
| Technical Sync | Weekly | Tech Lead + Backend | Architecture decisions |
| **Security Review** | **Weekly (Phase 3-4)** | **Tech Lead + Security** | **Security decisions** |

### 11.2 Reporting

| Report | Frequency | Audience |
|--------|-----------|----------|
| Sprint Status | Weekly | Stakeholders |
| Burndown Chart | Daily | Team |
| Risk Register | Bi-weekly | Management |
| Quality Metrics | Per sprint | Team + stakeholders |
| **Security Status** | **Weekly (Phase 3-4)** | **Management** |
| **Compliance Status** | **Bi-weekly (Phase 3-4)** | **Management** |

---

## 12. Acceptance Criteria Summary

### 12.1 Feature Acceptance

| Feature | Acceptance Criteria |
|---------|---------------------|
| Authentication | Login < 500ms, lockout after 5 failures, token blacklist works |
| Registration | Email verification, approval workflow, Masterdata validation |
| User Management | CRUD, status lifecycle, profile updates, password history |
| RBAC | Role assignment, permission check < 100ms, branch scoping |
| Organization | Multi-tenant isolation, branch scoping |
| Audit | All events logged, queryable via API, retention policy applied |

### 12.2 Non-Functional Acceptance

| Category | Criteria |
|----------|----------|
| Performance | p95 < 500ms, 1000 concurrent users |
| Security | 0 critical/high vulnerabilities, pen test passed |
| Availability | 99.9% uptime design |
| Documentation | API docs, integration guide, runbook |
| **Compliance** | ISO 27001 checklist passed, PCI DSS (if applicable) |

### 12.3 API Contract Acceptance

| Criteria | Measurement |
|----------|-------------|
| OpenAPI spec matches implementation | Postman tests 100% pass |
| All endpoints documented | OpenAPI coverage 100% |
| Error responses documented | All error codes in spec |
| Example requests/responses | All endpoints have examples |

---

## 13. Compliance Requirements

### 13.1 ISO 27001 Alignment

| Control Area | Requirement | Implementation |
|--------------|-------------|----------------|
| A.9.2.1 | User registration and de-registration | Approval workflow, status lifecycle |
| A.9.2.3 | Management of privileged access | RBAC, role assignment audit |
| A.9.4.1 | Information access restriction | Permission check, branch scoping |
| A.9.4.2 | Secure log-on procedures | MFA (PIN), rate limiting, lockout |
| A.12.4.1 | Event logging | Audit logging to OpenSearch |
| A.12.4.2 | Protection of log information | Immutable logs, access control |

### 13.2 PCI DSS Alignment (if applicable)

| Requirement | Implementation |
|-------------|----------------|
| 8.1.1 | Unique user IDs | UUID per user |
| 8.1.5 | Manage terminated user IDs | Status = INACTIVE, soft delete |
| 8.2.3 | Strong authentication | Password policy, PIN, MFA |
| 8.2.4 | Password complexity | Configurable per tenant |
| 10.1 | Audit trails for user access | All auth events logged |
| 10.2.1 | Log all individual access | Per-request audit logging |

> **See:** IAM-Compliance-Checklist.md for detailed compliance verification checklist.

---

## Document Sign-Off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Project Manager | | | |
| Tech Lead | | | |
| Product Owner | | | |
| Engineering Manager | | | |
| **Security Lead** | | | |

---

**End of Implementation Plan v2.0**

---

## Appendix A: Database Migration Files (16 files)

Per TRD-Section3-DatabaseDesign.md:

```
migrations/
├── 000001_create_trigger_function.up.sql
├── 000001_create_trigger_function.down.sql
├── 000002_create_tenants.up.sql
├── 000002_create_tenants.down.sql
├── 000003_create_branches.up.sql
├── 000003_create_branches.down.sql
├── 000004_create_users.up.sql
├── 000004_create_users.down.sql
├── 000005_create_user_auth_methods.up.sql
├── 000005_create_user_auth_methods.down.sql
├── 000006_create_user_profiles.up.sql
├── 000006_create_user_profiles.down.sql
├── 000007_create_user_security_states.up.sql
├── 000007_create_user_security_states.down.sql
├── 000008_create_user_branches.up.sql
├── 000008_create_user_branches.down.sql
├── 000009_create_password_history.up.sql
├── 000009_create_password_history.down.sql
├── 000010_create_applications.up.sql
├── 000010_create_applications.down.sql
├── 000011_create_roles.up.sql
├── 000011_create_roles.down.sql
├── 000012_create_permissions.up.sql
├── 000012_create_permissions.down.sql
├── 000013_create_role_permissions.up.sql
├── 000013_create_role_permissions.down.sql
├── 000014_create_user_role_assignments.up.sql
├── 000014_create_user_role_assignments.down.sql
├── 000015_create_indexes.up.sql
├── 000015_create_indexes.down.sql
├── 000016_seed_platform_data.up.sql
└── 000016_seed_platform_data.down.sql
```

---

## Appendix B: API Contract Deliverables

| Deliverable | Format | Location |
|-------------|--------|----------|
| OpenAPI Specification | YAML | openapi-iam-v1.yaml |
| Postman Collection | JSON | IAM-API-v1.postman_collection.json |
| Postman Environment (Dev) | JSON | IAM-Dev.postman_environment.json |
| Postman Environment (Staging) | JSON | IAM-Staging.postman_environment.json |
| Postman Environment (Prod) | JSON | IAM-Prod.postman_environment.json |

---

## Appendix C: Related Documents

| Document | Description |
|----------|-------------|
| BRD-GENERIC-IAM-SYSTEM-V1.md | Business Requirements Document |
| PRD-IAM-SYSTEM-V1.md | Product Requirements Document |
| PRD-Section2-F1-Authentication.md | Authentication feature spec |
| PRD-Section2-F2-UserManagement.md | User management feature spec |
| PRD-Section2-F3-Authorization.md | Authorization (RBAC) feature spec |
| PRD-Section2-F4-OrganizationManagement.md | Organization management feature spec |
| PRD-Section2-F5-AuditLogging.md | Audit logging feature spec |
| TRD-IAM-SYSTEM-V1.md | Technical Requirements Document |
| TRD-Section3-DatabaseDesign.md | Database design specification |
| TRD-Section4-APIDesign.md | API design specification |
| TRD-Section5-Security.md | Security architecture |
| IAM-ERD-Diagram.md | Entity Relationship Diagram |
| openapi-iam-v1.yaml | OpenAPI specification |
| **IAM-Compliance-Checklist.md** | Compliance verification checklist |
| **IAM-Detailed-Task-Breakdown.md** | Detailed task breakdown |
