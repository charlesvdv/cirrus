import Nav from './Nav'

export default (props) => {
  return (
    <>
      <Nav />
      <main>
        {props.children}
      </main>
    </>
  )
}