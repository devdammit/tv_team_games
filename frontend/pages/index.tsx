import Head from 'next/head'
import styles from '../styles/Home.module.css'
import {Avatar, Card, Skeleton, Tag, Typography} from "antd";
import {PointStatus, Result, useWatchPoints} from "../hooks/point";
import {useEffect, useState} from "react";

const positions = [...new Array(10)]

type GameStatus = {
    color?: string
    text: string
}

export default function Home() {
    const {status, results, startedAt, roundNum} = useWatchPoints()
    const [gameStatus, setGameStatus] = useState<GameStatus>({ text: "Ожидание" })


    useEffect(() => {
        if (status === PointStatus.STARTED) {
            setGameStatus({ color: 'green', text: 'В процессе' })
        } else if (status === PointStatus.FINISHED) {
            setGameStatus({ color: 'purple', text: 'Окончен' })
        }
    }, [status])

    return (
        <div className={styles.container}>
            <Head>
                <title>Victo UI</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <div className={styles.wrapper}>
                {positions.map((_, index) => {
                    const isWaiting =  !results || typeof results[index] === 'undefined'
                    const result = (!!results && results[index] ? results[index] : {}) as Result
                    let diffMs = 0
                    let color = 'green'

                    if (index > 3 && index < 7) {
                        color = 'yellow'
                    } else if (index >= 7) {
                        color = 'purple'
                    }
                        // `${diffMs} мс`

                    if (result && startedAt) {
                        const time1 = new Date(result.timestamp).getTime()
                        const time2 = startedAt.getTime()

                        diffMs = Math.round(time1 - time2)
                    }


                    return (
                        <div key={index} className={styles.team}>
                            <Card title={`${index + 1} место`} className={styles.teamCard}>
                                <Skeleton loading={isWaiting} avatar active>
                                    <Card.Meta
                                        avatar={
                                            <Avatar style={{ background: '#f56a00' }} >{result.name ? result.name[0] : ''}</Avatar>
                                        }
                                        title={result.name}
                                        description={<Tag color={color}>{`${diffMs} мс`}</Tag>}
                                    />
                                </Skeleton>
                            </Card>
                        </div>
                    )
                })}
            </div>

            <div className={styles.stateBar}>
                <div>
                    <Typography.Text>Состояние раунда: </Typography.Text>
                    <Typography.Text>
                        <Tag color={gameStatus.color}>{gameStatus.text}</Tag>
                    </Typography.Text>
                </div>
                <div>
                    <Typography.Text>Раунд №:</Typography.Text>
                    <Typography.Text keyboard>{roundNum + 1}</Typography.Text>
                </div>
            </div>
        </div>
    )
}
