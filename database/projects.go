package database

type Project struct {
	ID          int
	Name        string
	TechStack   string
	Deployed    string
	Description string
	GithubUrl   string
	Url 	    string
	Year 	    string
}

func GetAllProjects() ([]Project, error) {
	rows, err := DB.Query("SELECT id, name, tech_stack, deployed, description, github_url, url, year FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.TechStack, &p.Deployed, &p.Description, &p.GithubUrl, &p.Url, &p.Year); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func SearchProjects(query string) ([]Project, error) {
	sql := `SELECT id, name, tech_stack, deployed, description, github_url, url 
	        FROM projects 
	        WHERE name LIKE ? OR tech_stack LIKE ? OR description LIKE ?`

	rows, err := DB.Query(sql, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.TechStack, &p.Deployed, &p.Description, &p.GithubUrl, &p.Url, &p.Year)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
