package storage

import "labops/model"

func (s *Store) SaveUser(u model.User) error { return s.put("users", u.ID, u) }
func (s *Store) GetUser(id string) (*model.User, error) {
	var u model.User
	if err := s.get("users", id, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (s *Store) SaveEquipment(e model.Equipment) error { return s.put("equipment", e.ID, e) }
func (s *Store) GetEquipment(id string) (*model.Equipment, error) {
	var e model.Equipment
	if err := s.get("equipment", id, &e); err != nil {
		return nil, err
	}
	return &e, nil
}
