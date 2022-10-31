package b

// BlueprintOneToOneRoom contains a homeserver for multiroom feature tests.
var BlueprintMultiRoom = MustValidate(Blueprint{
	Name: "multiroom",
	Homeservers: []Homeserver{
		{
			Name: "hs1",
			Users: []User{
				{
					Localpart:   "@alice",
					DisplayName: "Alice",
				},
				{
					Localpart:   "@bob",
					DisplayName: "Bob",
				},
			},
		},
	},
})
