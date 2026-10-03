package blur_postgres_repository

import core_postgres_pool "github.com/vielchelpykh/facelab/internal/core/repository/postgres"

type BlurRepository struct {
	Pool *core_postgres_pool.Pool
}

func NewBlurRepository(pool *core_postgres_pool.Pool) *BlurRepository {
	return &BlurRepository{
		Pool: pool,
	}
}
