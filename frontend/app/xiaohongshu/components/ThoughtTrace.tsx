import styles from './thoughtTrace.module.css';

type ThoughtProps = {
  loading: boolean;
};

export default function ThoughtTrace(props: ThoughtProps) {
  // 旋转动画
  return (
    <div
      className={`${styles.box} ${props.loading ? styles.rotating : ''}`}
    />
  );
}
