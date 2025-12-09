-- +goose up
create table users (
	id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name VARCHAR(50) NOT NULL
);

-- +goose Down
drop table users;
