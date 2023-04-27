USE cambio_dolar;     

CREATE TABLE IF NOT EXISTS `cotacoes` (
`id` varchar(255),
`code` varchar(255) DEFAULT NULL,
`codein` varchar(50) DEFAULT NULL,
`name` varchar(50) DEFAULT NULL,
`bid` varchar(20) DEFAULT NULL
);