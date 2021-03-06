package manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/armakuni/circleci-context-secret-manager/manager"
)

var _ = Describe("Contexts", func() {
	Describe("HasAWSRemoteSecrets", func() {
		Context("when remote secrets are not enabled", func() {
			It("returns false", func() {
				var contexts = manager.Contexts{
					"main.yml": manager.Context{},
				}
				remoteSecretsEnabled, err := contexts.HasAWSRemoteSecrets()
				Ω(remoteSecretsEnabled).Should(BeFalse())
				Ω(err).Should(BeNil())
			})
		})

		Context("when remote secrets are enbaled", func() {
			Context("and the secret manager type is supported", func() {
				It("returns true", func() {
					var contexts = manager.Contexts{
						"main.yml": manager.Context{
							RemoteSecretStore: &manager.RemoteSecretStore{
								Type: "aws-secret-manager",
							},
						},
					}
					remoteSecretsEnabled, err := contexts.HasAWSRemoteSecrets()
					Ω(remoteSecretsEnabled).Should(BeTrue())
					Ω(err).Should(BeNil())
				})
			})

			Context("and the secret manager type is not supported", func() {
				It("returns an error", func() {
					var contexts = manager.Contexts{
						"main.yml": manager.Context{
							RemoteSecretStore: &manager.RemoteSecretStore{
								Type: "not-implemented",
							},
						},
					}
					remoteSecretsEnabled, err := contexts.HasAWSRemoteSecrets()
					Ω(remoteSecretsEnabled).Should(BeTrue())
					Ω(err).Should(MatchError("Unsupported remote secret manager 'not-implemented', supported managers are"))
				})
			})
		})
	})

	Describe("Process", func() {
		Context("when nothing extends", func() {
			var contexts = manager.Contexts{
				"main.yml": manager.Context{
					ContextID: "1",
					Secrets: manager.Secrets{
						"FOO": "bar",
						"BAZ": "fish",
					},
				},
				"dev.yml": manager.Context{
					ContextID: "2",
					Secrets: manager.Secrets{
						"TEST": "test",
					},
				},
			}

			It("returns contexts as is", func() {
				Ω(contexts.Process()).Should(Equal(contexts))
			})
		})

		Context("when a context extends another", func() {
			Context("and it has no overlapping secrets", func() {
				var contexts = manager.Contexts{
					"main.yml": manager.Context{
						ContextID: "1",
						Secrets: manager.Secrets{
							"FOO": "bar",
							"BAZ": "fish",
						},
					},
					"dev.yml": manager.Context{
						ContextID: "2",
						Extends:   []string{"main.yml"},
						Secrets: manager.Secrets{
							"TEST": "test",
						},
					},
					"test.yml": manager.Context{
						ContextID: "3",
						Extends:   []string{"main.yml", "dev.yml"},
						Secrets: manager.Secrets{
							"FOOBAR": "baz",
						},
					},
				}

				It("proccesses the extended contexts", func() {
					Ω(contexts.Process()).Should(Equal(manager.Contexts{
						"main.yml": manager.Context{
							ContextID: "1",
							Secrets: manager.Secrets{
								"FOO": "bar",
								"BAZ": "fish",
							},
						},
						"dev.yml": manager.Context{
							ContextID: "2",
							Secrets: manager.Secrets{
								"FOO":  "bar",
								"BAZ":  "fish",
								"TEST": "test",
							},
						},
						"test.yml": manager.Context{
							ContextID: "3",
							Secrets: manager.Secrets{
								"FOO":    "bar",
								"BAZ":    "fish",
								"TEST":   "test",
								"FOOBAR": "baz",
							},
						},
					}))
				})
			})

			Context("and it has overlapping secrets", func() {
				var contexts = manager.Contexts{
					"main.yml": manager.Context{
						ContextID: "1",
						Secrets: manager.Secrets{
							"FOO": "bar",
							"BAZ": "fish",
						},
					},
					"dev.yml": manager.Context{
						ContextID: "2",
						Extends:   []string{"main.yml"},
						Secrets: manager.Secrets{
							"TEST": "test",
							"BAZ":  "fwibble",
						},
					},
					"test.yml": manager.Context{
						ContextID: "3",
						Extends:   []string{"main.yml", "dev.yml"},
						Secrets: manager.Secrets{
							"FOOBAR": "baz",
							"FOO":    "jerry",
						},
					},
				}

				It("proccesses the extended contexts, overriding secrets where appropriate", func() {
					Ω(contexts.Process()).Should(Equal(manager.Contexts{
						"main.yml": manager.Context{
							ContextID: "1",
							Secrets: manager.Secrets{
								"FOO": "bar",
								"BAZ": "fish",
							},
						},
						"dev.yml": manager.Context{
							ContextID: "2",
							Secrets: manager.Secrets{
								"FOO":  "bar",
								"BAZ":  "fwibble",
								"TEST": "test",
							},
						},
						"test.yml": manager.Context{
							ContextID: "3",
							Secrets: manager.Secrets{
								"FOO":    "jerry",
								"BAZ":    "fwibble",
								"TEST":   "test",
								"FOOBAR": "baz",
							},
						},
					}))
				})
			})
		})
	})
})

var _ = Describe("Projects", func() {
	Describe("Process", func() {
		Context("when nothing extends", func() {
			var projects = manager.Projects{
				"main.yml": manager.Project{
					ProjectSlug: "github/org/repo1",
					Secrets: manager.Secrets{
						"FOO": "bar",
						"BAZ": "fish",
					},
				},
				"dev.yml": manager.Project{
					ProjectSlug: "github/org/repo2",
					Secrets: manager.Secrets{
						"TEST": "test",
					},
				},
			}

			It("returns contexts as is", func() {
				Ω(projects.Process()).Should(Equal(projects))
			})
		})

		Context("when a project extends another", func() {
			Context("and it has no overlapping secrets", func() {
				var projects = manager.Projects{
					"main.yml": manager.Project{
						ProjectSlug: "github/org/repo1",
						Secrets: manager.Secrets{
							"FOO": "bar",
							"BAZ": "fish",
						},
					},
					"dev.yml": manager.Project{
						ProjectSlug: "github/org/repo2",
						Extends:     []string{"main.yml"},
						Secrets: manager.Secrets{
							"TEST": "test",
						},
					},
					"test.yml": manager.Project{
						ProjectSlug: "github/org/repo3",
						Extends:     []string{"main.yml", "dev.yml"},
						Secrets: manager.Secrets{
							"FOOBAR": "baz",
						},
					},
				}

				It("proccesses the extended contexts", func() {
					Ω(projects.Process()).Should(Equal(manager.Projects{
						"main.yml": manager.Project{
							ProjectSlug: "github/org/repo1",
							Secrets: manager.Secrets{
								"FOO": "bar",
								"BAZ": "fish",
							},
						},
						"dev.yml": manager.Project{
							ProjectSlug: "github/org/repo2",
							Secrets: manager.Secrets{
								"FOO":  "bar",
								"BAZ":  "fish",
								"TEST": "test",
							},
						},
						"test.yml": manager.Project{
							ProjectSlug: "github/org/repo3",
							Secrets: manager.Secrets{
								"FOO":    "bar",
								"BAZ":    "fish",
								"TEST":   "test",
								"FOOBAR": "baz",
							},
						},
					}))
				})
			})

			Context("and it has overlapping secrets", func() {
				var projects = manager.Projects{
					"main.yml": manager.Project{
						ProjectSlug: "github/org/repo1",
						Secrets: manager.Secrets{
							"FOO": "bar",
							"BAZ": "fish",
						},
					},
					"dev.yml": manager.Project{
						ProjectSlug: "github/org/repo2",
						Extends:     []string{"main.yml"},
						Secrets: manager.Secrets{
							"TEST": "test",
							"BAZ":  "fwibble",
						},
					},
					"test.yml": manager.Project{
						ProjectSlug: "github/org/repo3",
						Extends:     []string{"main.yml", "dev.yml"},
						Secrets: manager.Secrets{
							"FOOBAR": "baz",
							"FOO":    "jerry",
						},
					},
				}

				It("proccesses the extended projects, overriding secrets where appropriate", func() {
					Ω(projects.Process()).Should(Equal(manager.Projects{
						"main.yml": manager.Project{
							ProjectSlug: "github/org/repo1",
							Secrets: manager.Secrets{
								"FOO": "bar",
								"BAZ": "fish",
							},
						},
						"dev.yml": manager.Project{
							ProjectSlug: "github/org/repo2",
							Secrets: manager.Secrets{
								"FOO":  "bar",
								"BAZ":  "fwibble",
								"TEST": "test",
							},
						},
						"test.yml": manager.Project{
							ProjectSlug: "github/org/repo3",
							Secrets: manager.Secrets{
								"FOO":    "jerry",
								"BAZ":    "fwibble",
								"TEST":   "test",
								"FOOBAR": "baz",
							},
						},
					}))
				})
			})
		})
	})
})
