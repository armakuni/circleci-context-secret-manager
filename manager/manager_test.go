package manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/armakuni/circleci-context-secret-manager/manager"
)

var _ = Describe("Contexts", func() {
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
