package shortener

type MemoryLinkStore struct {
	Links map[Identifier]Link
}

func (l MemoryLinkStore) GetLink(id Identifier, ip_address string, user_id int) (Link, error) {
	link, ok := l.Links[id]

	if !ok {
		return "", ErrLinkNotFound
	}

	return link, nil
}

func (l MemoryLinkStore) SaveLink(data *LinkPost) error {
	_, ok := l.Links[data.ID]

	if ok {
		return ErrIDExists
	}

	l.Links[data.ID] = data.Target

	return nil
}

func (l MemoryLinkStore) GetLinks(limit, offset int, orderBy string) error {
	return nil
}

func (l MemoryLinkStore) DeleteLink(id Identifier) error {
	delete(l.Links, id)

	return nil
}
