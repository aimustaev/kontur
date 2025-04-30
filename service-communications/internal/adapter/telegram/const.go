package telegram

func (a *TelegramAdapter) getProblems() map[int64][]Problem {
	problems := Problem{
		id:    0,
		label: "root",
		children: []Problem{
			{
				id:    1,
				label: "Тематика 1",
				children: []Problem{
					{
						id:       11,
						label:    "Тематика 11",
						children: []Problem{},
					},
					{
						id:       12,
						label:    "Тематика 12",
						children: []Problem{},
					},
				},
			},
			{
				id:       2,
				label:    "Тематика 2",
				children: []Problem{},
			},
			{
				id:       3,
				label:    "Тематика 3",
				children: []Problem{},
			},
		},
	}

	dictProblems := make(map[int64][]Problem)
	a.mapProblems(problems, &dictProblems)

	return dictProblems
}

func (a *TelegramAdapter) mapProblems(problem Problem, dict *map[int64][]Problem) {
	(*dict)[problem.id] = problem.children

	for _, pb := range problem.children {
		a.mapProblems(pb, dict)
	}
}
